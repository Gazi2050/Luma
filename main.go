package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ===== Token Types =====
type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	STRING = "STRING"
	BOOL   = "BOOL"

	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	ASSIGN   = "="

	LET    = "LET"
	CONST  = "CONST"
	LOG    = "LOG"
	FN     = "FN"
	LOOP   = "LOOP"
	RETURN = "RETURN"

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="
	EQ  = "=="
	NEQ = "!="

	SEMICOLON = ";"
	COMMA     = ","
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	COLON     = ":"
	LBRACKET  = "["
	RBRACKET  = "]"
)

// ===== Token =====
type Token struct {
	Type    TokenType
	Literal string
}

// ===== Lexer =====
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
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
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

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
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
	l.readChar()
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	str := l.input[start:l.position]
	l.readChar()
	return str
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func LookupIdent(ident string) TokenType {
	switch ident {
	case "let":
		return LET
	case "const":
		return CONST
	case "log":
		return LOG
	case "fn":
		return FN
	case "loop":
		return LOOP
	case "true", "false":
		return BOOL
	case "return":
		return RETURN
	default:
		return IDENT
	}
}

// ===== AST =====
type Statement interface{}
type Expression interface{}

type LetStatement struct {
	Name  string
	Value Expression
}

type LogStatement struct {
	Value Expression
}

type NumberLiteral struct {
	Value int
}

type StringLiteral struct {
	Value string
}

type BooleanLiteral struct {
	Value bool
}

type Identifier struct {
	Value string
}

type BinaryExpression struct {
	Left     Expression
	Operator TokenType
	Right    Expression
}

type ArrayLiteral struct {
	Elements []Expression
}

type ObjectLiteral struct {
	Pairs map[string]Expression
}

// ===== Environment =====
type Env struct {
	store map[string]interface{}
}

func NewEnv() *Env {
	return &Env{store: make(map[string]interface{})}
}

// ===== Evaluator =====
func Eval(stmt Statement, env *Env) interface{} {
	switch s := stmt.(type) {
	case *LetStatement:
		val := evalExpr(s.Value, env)
		env.store[s.Name] = val
		return val
	case *LogStatement:
		val := evalExpr(s.Value, env)
		fmt.Println(val)
		return val
	default:
		return nil
	}
}

func evalExpr(expr Expression, env *Env) interface{} {
	switch e := expr.(type) {
	case *NumberLiteral:
		return e.Value
	case *StringLiteral:
		return e.Value
	case *BooleanLiteral:
		return e.Value
	case *Identifier:
		return env.store[e.Value]
	case *BinaryExpression:
		left := evalExpr(e.Left, env).(int)
		right := evalExpr(e.Right, env).(int)
		switch e.Operator {
		case PLUS:
			return left + right
		case MINUS:
			return left - right
		case ASTERISK:
			return left * right
		case SLASH:
			return left / right
		}
	case *ArrayLiteral:
		var arr []interface{}
		for _, el := range e.Elements {
			arr = append(arr, evalExpr(el, env))
		}
		return arr
	case *ObjectLiteral:
		obj := make(map[string]interface{})
		for k, v := range e.Pairs {
			obj[k] = evalExpr(v, env)
		}
		return obj
	}
	return nil
}

// ===== Minimal Parser =====
func parse(tokens []Token) []Statement {
	var stmts []Statement
	i := 0
	for i < len(tokens) {
		tok := tokens[i]

		// let / const
		if tok.Type == LET || tok.Type == CONST {
			name := tokens[i+1].Literal
			i += 3
			var value Expression
			switch tokens[i].Type {
			case NUMBER:
				value = &NumberLiteral{Value: atoi(tokens[i].Literal)}
			case STRING:
				value = &StringLiteral{Value: tokens[i].Literal}
			case BOOL:
				value = &BooleanLiteral{Value: tokens[i].Literal == "true"}
			default:
				value = &Identifier{Value: tokens[i].Literal}
			}
			i++
			if i < len(tokens) && tokens[i].Type == SEMICOLON {
				i++
			}
			stmts = append(stmts, &LetStatement{Name: name, Value: value})

		} else if tok.Type == LOG {
			i++ // skip 'log'
			if i >= len(tokens) || tokens[i].Type != LPAREN {
				fmt.Println("Syntax error: expected '(' after log")
				break
			}
			i++ // skip '('

			// collect expression (simplified: single token expressions)
			exprTok := tokens[i]
			var expr Expression
			switch exprTok.Type {
			case NUMBER:
				expr = &NumberLiteral{Value: atoi(exprTok.Literal)}
			case STRING:
				expr = &StringLiteral{Value: exprTok.Literal}
			case BOOL:
				expr = &BooleanLiteral{Value: exprTok.Literal == "true"}
			case IDENT:
				expr = &Identifier{Value: exprTok.Literal}
			default:
				fmt.Println("Unsupported log argument")
				i++
				continue
			}
			i++ // move past expression
			if i < len(tokens) && tokens[i].Type == RPAREN {
				i++
			}
			if i < len(tokens) && tokens[i].Type == SEMICOLON {
				i++
			}
			stmts = append(stmts, &LogStatement{Value: expr})

		} else {
			i++
		}
	}
	return stmts
}

func atoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

// ===== REPL =====
func main() {
	env := NewEnv()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Luma REPL v1 - Type your code")
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		lexer := NewLexer(line)
		var tokens []Token
		for tok := lexer.NextToken(); tok.Type != EOF; tok = lexer.NextToken() {
			tokens = append(tokens, tok)
		}
		stmts := parse(tokens)
		for _, stmt := range stmts {
			Eval(stmt, env)
		}
	}
}
