package parser

import (
	"luma/ast"
	"luma/lexer"
	"luma/token"
	"strconv"
	"strings"
)

const (
	_ int = iota
	LOWEST
	ASSIGNMENT  // =
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.DOT:      INDEX,
	token.ASSIGN:   ASSIGNMENT,
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET, token.CONST:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.LOOP:
		return p.parseLoopStatement()
	case token.FN:
		if p.peekToken.Type == token.IDENT {
			return p.parseFunctionStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFunctionStatement() ast.Statement {
	p.nextToken() // fn
	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	lit := &ast.FunctionLiteral{Token: token.Token{Type: token.FN, Literal: "fn"}}
	if !p.expectPeek(token.LPAREN) { return nil }
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) { return nil }
	lit.Body = p.parseBlockStatement()

	return &ast.LetStatement{
		Token: token.Token{Type: token.LET, Literal: "let"},
		Name:  name,
		Value: lit,
		IsConst: false,
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken}
	stmt.IsConst = p.curToken.Type == token.CONST

	if p.peekToken.Type != token.IDENT {
		return nil
	}
	p.nextToken()

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != token.ASSIGN {
		return nil
	}
	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil { return nil }

	if p.peekToken.Type == token.SEMICOLON { p.nextToken() }

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil { return nil }
	if p.peekToken.Type == token.SEMICOLON { p.nextToken() }
	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if stmt.Expression == nil { return nil }
	if p.peekToken.Type == token.SEMICOLON { p.nextToken() }
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	p.nextToken()
	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseLoopStatement() ast.Statement {
	stmt := &ast.LoopStatement{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) { return nil }
	p.nextToken() // init

	stmt.Init = p.parseStatement()
	if stmt.Init == nil { return nil }

	if p.curToken.Type != token.SEMICOLON {
		if !p.expectPeek(token.SEMICOLON) { return nil }
	}
	p.nextToken() // skip ;

	stmt.Condition = p.parseExpression(LOWEST)
	if stmt.Condition == nil { return nil }
	if !p.expectPeek(token.SEMICOLON) { return nil }
	p.nextToken() // skip ;

	stmt.Post = p.parseExpression(LOWEST)
	if stmt.Post == nil { return nil }
	if !p.expectPeek(token.RPAREN) { return nil }
	
	if !p.expectPeek(token.LBRACE) { return nil }
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var leftExp ast.Expression
	
	// Prefix
	switch p.curToken.Type {
	case token.IDENT:
		leftExp = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case token.NUMBER:
		v, _ := strconv.Atoi(p.curToken.Literal)
		leftExp = &ast.NumberLiteral{Token: p.curToken, Value: v}
	case token.STRING:
		leftExp = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case token.TRUE, token.FALSE:
		leftExp = &ast.BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == token.TRUE}
	case token.MINUS:
		leftExp = p.parsePrefixExpression()
	case token.LPAREN:
		p.nextToken()
		exp := p.parseExpression(LOWEST)
		if p.peekToken.Type != token.RPAREN { return nil }
		p.nextToken()
		leftExp = exp
	case token.LBRACKET:
		leftExp = p.parseArrayLiteral()
	case token.LBRACE:
		leftExp = p.parseObjectLiteral()
	case token.FN:
		leftExp = p.parseFunctionLiteral()
	case token.IF:
		leftExp = p.parseIfExpression()
	case token.LOG, token.LEN:
        leftExp = &ast.Identifier{Token: p.curToken, Value: strings.ToLower(p.curToken.Literal)}
	}

	for p.peekToken.Type != token.SEMICOLON && precedence < precedences[p.peekToken.Type] {
		p.nextToken()
		leftExp = p.parseInfixExpression(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	switch p.curToken.Type {
	case token.LPAREN:
		return p.parseCallExpression(left)
	case token.LBRACKET:
		exp := &ast.IndexExpression{Token: p.curToken, Left: left}
		p.nextToken()
		exp.Index = p.parseExpression(LOWEST)
		if p.peekToken.Type != token.RBRACKET { return nil }
		p.nextToken()
		return exp
	case token.DOT:
		exp := &ast.MemberExpression{Token: p.curToken, Left: left}
		p.nextToken()
		exp.Member = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		return exp
	default:
		exp := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}
		prec := precedences[p.curToken.Type]
		p.nextToken()
		exp.Right = p.parseExpression(prec)
		return exp
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	al := &ast.ArrayLiteral{Token: p.curToken}
	al.Elements = p.parseExpressionList(token.RBRACKET)
	return al
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.curToken, Pairs: make(map[string]ast.Expression)}
	for p.peekToken.Type != token.RBRACE {
		p.nextToken()
		key := p.curToken.Literal
		if p.peekToken.Type != token.COLON { return nil }
		p.nextToken() 
		p.nextToken()
		obj.Pairs[key] = p.parseExpression(LOWEST)
		if p.peekToken.Type != token.COMMA && p.peekToken.Type != token.RBRACE { return nil }
		if p.peekToken.Type == token.COMMA { p.nextToken() }
	}
	p.nextToken() 
	return obj
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekToken.Type == end {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if p.peekToken.Type != end { return nil }
	p.nextToken()
	return list
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if p.peekToken.Type != token.LPAREN { return nil }
	p.nextToken()
	lit.Parameters = p.parseFunctionParameters()
	if p.peekToken.Type != token.LBRACE { return nil }
	p.nextToken()
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}
	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return params
	}
	p.nextToken()
	params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}
	if p.peekToken.Type != token.RPAREN { return nil }
	p.nextToken()
	return params
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}
	if p.peekToken.Type != token.LPAREN { return nil }
	p.nextToken()
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)
	if p.peekToken.Type != token.RPAREN { return nil }
	p.nextToken()
	if p.peekToken.Type != token.LBRACE { return nil }
	p.nextToken()
	exp.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == token.ELSE {
		p.nextToken()
		if p.peekToken.Type != token.LBRACE { return nil }
		p.nextToken()
		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}
