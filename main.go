package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Token types
type TokenType string

const (
	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	STRING = "STRING"

	// Separators
	SEMICOLON = ";"
	COMMA     = ","
	COLON     = ":"
	DOT       = "."
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"

	// Comparison
	EQ  = "=="
	NEQ = "!="
	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	// Keywords
	LET    = "LET"
	CONST  = "CONST"
	FN     = "FN"
	RETURN = "RETURN"
	IF     = "IF"
	ELSE   = "ELSE"
	LOOP   = "LOOP"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	// Built-in placeholder (though we handle them as identifiers usually)
	LOG = "LOG"
	LEN = "LEN"

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

// Token
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' || (l.ch == '/' && l.peekChar() == '/') {
		if l.ch == '/' && l.peekChar() == '/' {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		} else {
			l.readChar()
		}
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readString() string {
	l.readChar() // skip "
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	str := l.input[start:l.position]
	if l.ch == '"' {
		l.readChar()
	}
	return str
}

func LookupIdent(ident string) TokenType {
	switch ident {
	case "let":
		return LET
	case "const":
		return CONST
	case "fn":
		return FN
	case "return":
		return RETURN
	case "if":
		return IF
	case "else":
		return ELSE
	case "loop":
		return LOOP
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "log":
		return LOG
	case "len":
		return LEN
	default:
		return IDENT
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()
	switch l.ch {
	case '+':
		tok = Token{Type: PLUS, Literal: string(l.ch)}
	case '-':
		tok = Token{Type: MINUS, Literal: string(l.ch)}
	case '*':
		tok = Token{Type: ASTERISK, Literal: string(l.ch)}
	case '/':
		tok = Token{Type: SLASH, Literal: string(l.ch)}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: EQ, Literal: "=="}
		} else {
			tok = Token{Type: ASSIGN, Literal: string(l.ch)}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: NEQ, Literal: "!="}
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch)}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: LTE, Literal: "<="}
		} else {
			tok = Token{Type: LT, Literal: "<"}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: GTE, Literal: ">="}
		} else {
			tok = Token{Type: GT, Literal: ">"}
		}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch)}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch)}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch)}
	case '[':
		tok = Token{Type: LBRACKET, Literal: string(l.ch)}
	case ']':
		tok = Token{Type: RBRACKET, Literal: string(l.ch)}
	case ',':
		tok = Token{Type: COMMA, Literal: string(l.ch)}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(l.ch)}
	case ':':
		tok = Token{Type: COLON, Literal: string(l.ch)}
	case '.':
		tok = Token{Type: DOT, Literal: string(l.ch)}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		return tok
	case 0:
		tok.Type = EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tok.Type = LookupIdent(lit)
			tok.Literal = lit
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch)}
		}
	}
	l.readChar()
	return tok
}

// --- AST ---

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token Token // LET or CONST
	Name  *Identifier
	Value Expression
	IsConst bool
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	kind := "let "
	if ls.IsConst { kind = "const " }
	return kind + ls.Name.String() + " = " + ls.Value.String() + ";"
}

type ReturnStatement struct {
	Token Token
	Value Expression
}
func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string { return "return " + rs.Value.String() + ";" }

type ExpressionStatement struct {
	Token Token
	Expression Expression
}
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string { return es.Expression.String() }

type BlockStatement struct {
	Token      Token
	Statements []Statement
}
func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	var out strings.Builder
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")
	return out.String()
}

type LoopStatement struct {
	Token     Token
	Init      Statement
	Condition Expression
	Post      Expression
	Body      *BlockStatement
}
func (ls *LoopStatement) statementNode() {}
func (ls *LoopStatement) String() string {
	return "loop (" + ls.Init.String() + " " + ls.Condition.String() + "; " + ls.Post.String() + ") " + ls.Body.String()
}

type Identifier struct {
	Token Token
	Value string
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string { return i.Value }

type NumberLiteral struct {
	Token Token
	Value int
}
func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string { return nl.Token.Literal }

type StringLiteral struct {
	Token Token
	Value string
}
func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string { return "\"" + sl.Value + "\"" }

type BooleanLiteral struct {
	Token Token
	Value bool
}
func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string { return bl.Token.Literal }

type ArrayLiteral struct {
	Token    Token
	Elements []Expression
}
func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) String() string {
	var out strings.Builder
	out.WriteString("[")
	for i, el := range al.Elements {
		out.WriteString(el.String())
		if i < len(al.Elements)-1 { out.WriteString(", ") }
	}
	out.WriteString("]")
	return out.String()
}

type ObjectLiteral struct {
	Token Token
	Pairs map[string]Expression
}
func (ol *ObjectLiteral) expressionNode() {}
func (ol *ObjectLiteral) String() string {
	var out strings.Builder
	out.WriteString("{")
	for k, v := range ol.Pairs {
		out.WriteString(k + ": " + v.String() + ", ")
	}
	out.WriteString("}")
	return out.String()
}

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string { return "(" + pe.Operator + pe.Right.String() + ")" }

type InfixExpression struct {
	Token    Token
	Left     Expression
	Operator string
	Right    Expression
}
func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string { return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")" }

type IfExpression struct {
	Token     Token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}
func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) String() string {
	out := "if " + ie.Condition.String() + " " + ie.Consequence.String()
	if ie.Alternative != nil {
		out += " else " + ie.Alternative.String()
	}
	return out
}

type FunctionLiteral struct {
	Token      Token
	Parameters []*Identifier
	Body       *BlockStatement
}
func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) String() string {
	var out strings.Builder
	out.WriteString("fn(")
	for i, p := range fl.Parameters {
		out.WriteString(p.String())
		if i < len(fl.Parameters)-1 { out.WriteString(", ") }
	}
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     Token
	Function  Expression // In Luma log() is handled via name
	Arguments []Expression
}
func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	var out strings.Builder
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for i, a := range ce.Arguments {
		out.WriteString(a.String())
		if i < len(ce.Arguments)-1 { out.WriteString(", ") }
	}
	out.WriteString(")")
	return out.String()
}

type IndexExpression struct {
	Token Token
	Left  Expression
	Index Expression
}
func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string { return "(" + ie.Left.String() + "[" + ie.Index.String() + "])" }

type MemberExpression struct {
	Token Token
	Left  Expression
	Member *Identifier
}
func (me *MemberExpression) expressionNode() {}
func (me *MemberExpression) String() string { return me.Left.String() + "." + me.Member.String() }

// --- Parser ---

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

var precedences = map[TokenType]int{
	EQ:       EQUALS,
	NEQ:      EQUALS,
	LT:       LESSGREATER,
	GT:       LESSGREATER,
	LTE:      LESSGREATER,
	GTE:      LESSGREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	ASTERISK: PRODUCT,
	LPAREN:   CALL,
	LBRACKET: INDEX,
	DOT:      INDEX, // member access is high precedence
	ASSIGN:   ASSIGNMENT,
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case LET, CONST:
		return p.parseLetStatement()
	case RETURN:
		return p.parseReturnStatement()
	case LOOP:
		return p.parseLoopStatement()
	case FN:
		if p.peekToken.Type == IDENT {
			return p.parseFunctionStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFunctionStatement() Statement {
	p.nextToken() // fn
	name := &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	lit := &FunctionLiteral{Token: Token{Type: FN, Literal: "fn"}}
	if !p.expectPeek(LPAREN) { return nil }
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(LBRACE) { return nil }
	lit.Body = p.parseBlockStatement()

	return &LetStatement{
		Token: Token{Type: LET, Literal: "let"},
		Name:  name,
		Value: lit,
		IsConst: false,
	}
}

func (p *Parser) parseLetStatement() Statement {
	stmt := &LetStatement{Token: p.curToken}
	stmt.IsConst = p.curToken.Type == CONST

	if p.peekToken.Type != IDENT {
		fmt.Printf("PARSER ERROR: expected IDENT after let/const, got %s\n", p.peekToken.Type)
		return nil
	}
	p.nextToken()

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekToken.Type != ASSIGN {
		fmt.Printf("PARSER ERROR: expected = after identifier, got %s\n", p.peekToken.Type)
		return nil
	}
	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil { return nil }

	if p.peekToken.Type == SEMICOLON { p.nextToken() }

	return stmt
}

func (p *Parser) parseReturnStatement() Statement {
	stmt := &ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil { return nil }
	if p.peekToken.Type == SEMICOLON { p.nextToken() }
	return stmt
}

func (p *Parser) parseExpressionStatement() Statement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if stmt.Expression == nil { return nil }
	if p.peekToken.Type == SEMICOLON { p.nextToken() }
	return stmt
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	p.nextToken()
	for p.curToken.Type != RBRACE && p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseLoopStatement() Statement {
	stmt := &LoopStatement{Token: p.curToken}
	if !p.expectPeek(LPAREN) { return nil }
	p.nextToken() // init

	stmt.Init = p.parseStatement()
	if stmt.Init == nil { return nil }

	if p.curToken.Type != SEMICOLON {
		if !p.expectPeek(SEMICOLON) { return nil }
	}
	p.nextToken() // skip ;

	stmt.Condition = p.parseExpression(LOWEST)
	if stmt.Condition == nil { return nil }
	if !p.expectPeek(SEMICOLON) { return nil }
	p.nextToken() // skip ;

	stmt.Post = p.parseExpression(LOWEST)
	if stmt.Post == nil { return nil }
	if !p.expectPeek(RPAREN) { return nil }
	
	if !p.expectPeek(LBRACE) { return nil }
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	var leftExp Expression
	
	// Prefix
	switch p.curToken.Type {
	case IDENT:
		leftExp = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case NUMBER:
		v, _ := strconv.Atoi(p.curToken.Literal)
		leftExp = &NumberLiteral{Token: p.curToken, Value: v}
	case STRING:
		leftExp = &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case TRUE, FALSE:
		leftExp = &BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == TRUE}
	case MINUS:
		leftExp = p.parsePrefixExpression()
	case LPAREN:
		p.nextToken()
		exp := p.parseExpression(LOWEST)
		if p.peekToken.Type != RPAREN { return nil }
		p.nextToken()
		leftExp = exp
	case LBRACKET:
		leftExp = p.parseArrayLiteral()
	case LBRACE:
		leftExp = p.parseObjectLiteral()
	case FN:
		leftExp = p.parseFunctionLiteral()
	case IF:
		leftExp = p.parseIfExpression()
	case LOG, LEN:
        // Built-ins can be treated as identifiers for call expression
        leftExp = &Identifier{Token: p.curToken, Value: strings.ToLower(p.curToken.Literal)}
	}

	for p.peekToken.Type != SEMICOLON && precedence < precedences[p.peekToken.Type] {
		p.nextToken()
		leftExp = p.parseInfixExpression(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() Expression {
	exp := &PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	switch p.curToken.Type {
	case LPAREN:
		return p.parseCallExpression(left)
	case LBRACKET:
		exp := &IndexExpression{Token: p.curToken, Left: left}
		p.nextToken()
		exp.Index = p.parseExpression(LOWEST)
		if p.peekToken.Type != RBRACKET { return nil }
		p.nextToken()
		return exp
	case DOT:
		exp := &MemberExpression{Token: p.curToken, Left: left}
		p.nextToken()
		exp.Member = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		return exp
	default:
		exp := &InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}
		prec := precedences[p.curToken.Type]
		p.nextToken()
		exp.Right = p.parseExpression(prec)
		return exp
	}
}

func (p *Parser) parseArrayLiteral() Expression {
	al := &ArrayLiteral{Token: p.curToken}
	al.Elements = p.parseExpressionList(RBRACKET)
	return al
}

func (p *Parser) parseObjectLiteral() Expression {
	obj := &ObjectLiteral{Token: p.curToken, Pairs: make(map[string]Expression)}
	for p.peekToken.Type != RBRACE {
		p.nextToken()
		key := p.curToken.Literal
		if p.peekToken.Type != COLON { return nil }
		p.nextToken() // :
		p.nextToken()
		obj.Pairs[key] = p.parseExpression(LOWEST)
		if p.peekToken.Type != COMMA && p.peekToken.Type != RBRACE { return nil }
		if p.peekToken.Type == COMMA { p.nextToken() }
	}
	p.nextToken() // }
	return obj
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}
	if p.peekToken.Type == end {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if p.peekToken.Type != end { return nil }
	p.nextToken()
	return list
}

func (p *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: p.curToken}
	if p.peekToken.Type != LPAREN { return nil }
	p.nextToken()
	lit.Parameters = p.parseFunctionParameters()
	if p.peekToken.Type != LBRACE { return nil }
	p.nextToken()
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	params := []*Identifier{}
	if p.peekToken.Type == RPAREN {
		p.nextToken()
		return params
	}
	p.nextToken()
	params = append(params, &Identifier{Token: p.curToken, Value: p.curToken.Literal})
	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		params = append(params, &Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}
	if p.peekToken.Type != RPAREN { return nil }
	p.nextToken()
	return params
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseIfExpression() Expression {
	exp := &IfExpression{Token: p.curToken}
	if p.peekToken.Type != LPAREN { return nil }
	p.nextToken() // (
	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)
	if p.peekToken.Type != RPAREN { return nil }
	p.nextToken() // )
	if p.peekToken.Type != LBRACE { return nil }
	p.nextToken() // {
	exp.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == ELSE {
		p.nextToken()
		if p.peekToken.Type != LBRACE { return nil }
		p.nextToken()
		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

// --- Environment ---

type Env struct {
	store map[string]interface{}
	consts map[string]bool
	outer *Env
}

func NewEnv() *Env {
	return &Env{store: make(map[string]interface{}), consts: make(map[string]bool)}
}

func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

func (e *Env) Get(name string) (interface{}, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return val, ok
}

func (e *Env) Set(name string, val interface{}, isConst bool) {
	e.store[name] = val
	if isConst { e.consts[name] = true }
}

func (e *Env) Update(name string, val interface{}) bool {
	if _, ok := e.store[name]; ok {
		if _, isConst := e.consts[name]; isConst {
			fmt.Printf("ERROR: Cannot reassign constant %s\n", name)
			return false
		}
		e.store[name] = val
		return true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return false
}

// --- Evaluator ---

type ReturnValue struct{ Value interface{} }

func Eval(node Node, env *Env) interface{} {
	if node == nil { return nil }
	val := reflect.ValueOf(node)
	if val.Kind() == reflect.Ptr && val.IsNil() { return nil }
	
	switch n := node.(type) {
	case *Program:
		var result interface{}
		for _, stmt := range n.Statements {
			result = Eval(stmt, env)
			if rv, ok := result.(*ReturnValue); ok { return rv.Value }
		}
		return result
	case *BlockStatement:
		var result interface{}
		for _, stmt := range n.Statements {
			result = Eval(stmt, env)
			if result != nil {
				if _, ok := result.(*ReturnValue); ok { return result }
			}
		}
		return result
	case *ExpressionStatement:
		return Eval(n.Expression, env)
	case *ReturnStatement:
		val := Eval(n.Value, env)
		return &ReturnValue{Value: val}
	case *LetStatement:
		val := Eval(n.Value, env)
		env.Set(n.Name.Value, val, n.IsConst)
		return val
	case *LoopStatement:
		loopEnv := NewEnclosedEnv(env)
		Eval(n.Init, loopEnv)
		for {
			cond := Eval(n.Condition, loopEnv)
			condBool, ok := cond.(bool)
			if !ok || !condBool { break }
			Eval(n.Body, loopEnv)
			Eval(n.Post, loopEnv)
		}
		return nil

	case *NumberLiteral: return n.Value
	case *StringLiteral: return n.Value
	case *BooleanLiteral: return n.Value
	case *Identifier:
		val, ok := env.Get(n.Value)
		if !ok {
			// Built-ins check
			if n.Value == "log" { return "BUILTIN_LOG" }
			if n.Value == "len" { return "BUILTIN_LEN" }
			fmt.Printf("ERROR: identifier not found: %s\n", n.Value)
			return nil
		}
		return val
	case *ArrayLiteral:
		elements := []interface{}{}
		for _, el := range n.Elements { elements = append(elements, Eval(el, env)) }
		return elements
	case *ObjectLiteral:
		obj := make(map[string]interface{})
		for k, v := range n.Pairs { obj[k] = Eval(v, env) }
		return obj

	case *PrefixExpression:
		right := Eval(n.Right, env)
		if right == nil { return nil }
		if n.Operator == "-" { 
			if r, ok := right.(int); ok { return -r }
		}
		return nil
	case *InfixExpression:
		left := Eval(n.Left, env)
		right := Eval(n.Right, env)
		if left == nil || right == nil { return nil }
		switch n.Operator {
		case "+":
			if l, ok := left.(string); ok { 
				if r, ok := right.(string); ok { return l + r }
			}
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok { return l + r }
			}
			return nil
		case "-": 
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok { return l - r }
			}
			return nil
		case "*":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok { return l * r }
			}
			return nil
		case "/":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok { return l / r }
			}
			return nil
		case "=":
			if ident, ok := n.Left.(*Identifier); ok {
				val := Eval(n.Right, env)
				env.Update(ident.Value, val)
				return val
			}
			return nil
		case "==": return left == right
		case "!=": return left != right
		case "<": 
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok { return l < r }
			}
			return false
		case ">":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok { return l > r }
			}
			return false
		case "<=":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok { return l <= r }
			}
			return false
		case ">=":
			if l, ok := left.(int); ok {
				if r, ok := right.(int); ok { return l >= r }
			}
			return false
		}
	case *IfExpression:
		cond := Eval(n.Condition, env)
		if cond.(bool) {
			return Eval(n.Consequence, env)
		} else if n.Alternative != nil {
			return Eval(n.Alternative, env)
		}
		return nil
	case *FunctionLiteral:
		return n // Function is its own AST node as a value
	case *CallExpression:
		fn := Eval(n.Function, env)
		args := []interface{}{}
		for _, a := range n.Arguments { args = append(args, Eval(a, env)) }

		if s, ok := fn.(string); ok {
			if s == "BUILTIN_LOG" {
				fmt.Println(args...)
				return nil
			}
			if s == "BUILTIN_LEN" {
				if arr, ok := args[0].([]interface{}); ok { return len(arr) }
				if str, ok := args[0].(string); ok { return len(str) }
				return 0
			}
		}

		if fl, ok := fn.(*FunctionLiteral); ok {
			childEnv := NewEnclosedEnv(env)
			for i, p := range fl.Parameters {
				childEnv.Set(p.Value, args[i], false)
			}
			res := Eval(fl.Body, childEnv)
			if rv, ok := res.(*ReturnValue); ok { return rv.Value }
			return res
		}
	case *IndexExpression:
		left := Eval(n.Left, env)
		idx := Eval(n.Index, env)
		if arr, ok := left.([]interface{}); ok { return arr[idx.(int)] }
		return nil
	case *MemberExpression:
		left := Eval(n.Left, env)
		if obj, ok := left.(map[string]interface{}); ok {
			return obj[n.Member.Value]
		}
		return nil
	}
	return nil
}

// --- Main / REPL ---

func runLumaScript(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	l := NewLexer(string(content))
	p := NewParser(l)
	program := p.ParseProgram()
	env := NewEnv()
	Eval(program, env)
}

func main() {
	if len(os.Args) > 2 && os.Args[1] == "run" {
		runLumaScript(os.Args[2])
		return
	}

	fmt.Println("Luma REPL v1")
	env := NewEnv()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		if !scanner.Scan() { break }
		line := scanner.Text()
		if line == "" { continue }
		l := NewLexer(line)
		p := NewParser(l)
		program := p.ParseProgram()
		Eval(program, env)
	}
}
