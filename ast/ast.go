package ast

import (
	"luma/token"
	"strings"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	StatementNode()
}

type Expression interface {
	Node
	ExpressionNode()
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
	Token   token.Token // LET or CONST
	Name    *Identifier
	Value   Expression
	IsConst bool
}

func (ls *LetStatement) StatementNode() {}
func (ls *LetStatement) String() string {
	kind := "let "
	if ls.IsConst {
		kind = "const "
	}
	return kind + ls.Name.String() + " = " + ls.Value.String() + ";"
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) StatementNode() {}
func (rs *ReturnStatement) String() string { return "return " + rs.Value.String() + ";" }

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) StatementNode() {}
func (es *ExpressionStatement) String() string { return es.Expression.String() }

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) StatementNode() {}
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
	Token     token.Token
	Init      Statement
	Condition Expression
	Post      Expression
	Body      *BlockStatement
}

func (ls *LoopStatement) StatementNode() {}
func (ls *LoopStatement) String() string {
	return "loop (" + ls.Init.String() + " " + ls.Condition.String() + "; " + ls.Post.String() + ") " + ls.Body.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) ExpressionNode() {}
func (i *Identifier) String() string  { return i.Value }

type NumberLiteral struct {
	Token token.Token
	Value int
}

func (nl *NumberLiteral) ExpressionNode() {}
func (nl *NumberLiteral) String() string  { return nl.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) ExpressionNode() {}
func (sl *StringLiteral) String() string  { return "\"" + sl.Value + "\"" }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) ExpressionNode() {}
func (bl *BooleanLiteral) String() string  { return bl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) ExpressionNode() {}
func (al *ArrayLiteral) String() string {
	var out strings.Builder
	out.WriteString("[")
	for i, el := range al.Elements {
		out.WriteString(el.String())
		if i < len(al.Elements)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString("]")
	return out.String()
}

type ObjectLiteral struct {
	Token token.Token
	Pairs map[string]Expression
}

func (ol *ObjectLiteral) ExpressionNode() {}
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
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) ExpressionNode() {}
func (pe *PrefixExpression) String() string  { return "(" + pe.Operator + pe.Right.String() + ")" }

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) ExpressionNode() {}
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) ExpressionNode() {}
func (ie *IfExpression) String() string {
	out := "if " + ie.Condition.String() + " " + ie.Consequence.String()
	if ie.Alternative != nil {
		out += " else " + ie.Alternative.String()
	}
	return out
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) ExpressionNode() {}
func (fl *FunctionLiteral) String() string {
	var out strings.Builder
	out.WriteString("fn(")
	for i, p := range fl.Parameters {
		out.WriteString(p.String())
		if i < len(fl.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) ExpressionNode() {}
func (ce *CallExpression) String() string {
	var out strings.Builder
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for i, a := range ce.Arguments {
		out.WriteString(a.String())
		if i < len(ce.Arguments)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) ExpressionNode() {}
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}

type MemberExpression struct {
	Token  token.Token
	Left   Expression
	Member *Identifier
}

func (me *MemberExpression) ExpressionNode() {}
func (me *MemberExpression) String() string  { return me.Left.String() + "." + me.Member.String() }
