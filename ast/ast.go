package ast

import (
	"bytes"
	"strings"
	"yap/token"
)

type Node interface {
	String() string
	TokenLiteral() string
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

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	var msg bytes.Buffer
	msg.WriteString("(")
	msg.WriteString(p.Operator)
	msg.WriteString(p.Right.String())
	msg.WriteString(")")

	return msg.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}

func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	var msg bytes.Buffer
	msg.WriteString("(")
	msg.WriteString(i.Left.String())
	msg.WriteString(" " + i.Operator + " ")
	msg.WriteString(i.Right.String())
	msg.WriteString(")")

	return msg.String()
}

type PostfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
}

func (p *PostfixExpression) expressionNode() {}
func (p *PostfixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PostfixExpression) String() string {
	var msg bytes.Buffer
	msg.WriteString("(")
	msg.WriteString(p.Left.String())
	msg.WriteString(" " + p.Operator + " ")
	msg.WriteString(p.Operator)
	msg.WriteString(")")

	return msg.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type SayStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (s *SayStatement) statementNode() {}

func (s *SayStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *SayStatement) String() string {
	var out bytes.Buffer

	out.WriteString(s.TokenLiteral() + " ")
	out.WriteString(s.Name.String())
	out.WriteString(" = ")

	if s.Value != nil {
		out.WriteString(s.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type PotentialStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (p *PotentialStatement) statementNode() {}
func (p *PotentialStatement) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PotentialStatement) String() string {
	var msg bytes.Buffer

	msg.WriteString(p.TokenLiteral() + " ")
	msg.WriteString(p.Name.String())
	msg.WriteString(" = ")
	msg.WriteString(p.Value.String())
	msg.WriteString(";")

	return msg.String()
}

type ConstStaement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (c *ConstStaement) statementNode() {}

func (c *ConstStaement) TokenLiteral() string {
	return c.Token.Literal
}

func (c *ConstStaement) String() string {
	var out bytes.Buffer

	out.WriteString(c.TokenLiteral() + " ")
	out.WriteString(c.Name.String())
	out.WriteString(" = ")

	if c.Value != nil {
		out.WriteString(c.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type GlobalStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (g *GlobalStatement) statementNode() {}

func (g *GlobalStatement) TokenLiteral() string {
	return g.Token.Literal
}

func (g *GlobalStatement) String() string {
	var out bytes.Buffer

	out.WriteString(g.TokenLiteral() + " ")
	out.WriteString(g.Name.String())
	out.WriteString(" = ")

	if g.Value != nil {
		out.WriteString(g.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (f *FloatLiteral) expressionNode() {}
func (f *FloatLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FloatLiteral) String() string {
	return f.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Elif        []ElifEpxression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IfExpression) String() string {
	var msg bytes.Buffer

	msg.WriteString("if")
	msg.WriteString(i.Condition.String())
	msg.WriteString(" ")
	msg.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		msg.WriteString("else")
		msg.WriteString(i.Alternative.String())
	}
	return msg.String()
}

type ElifEpxression struct {
	Token        token.Token
	Conditions   Expression
	Consequences *BlockStatement
}

type TernaryExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (t *TernaryExpression) expressionNode() {}
func (t *TernaryExpression) TokenLiteral() string {
	return t.Token.Literal
}

func (t *TernaryExpression) String() string {
	var msg bytes.Buffer
	msg.WriteString("if")
	msg.WriteString(t.Condition.String())
	msg.WriteString("then")
	msg.WriteString(t.Consequence.String())

	if t.Alternative != nil {
		msg.WriteString("else")
		msg.WriteString(t.Alternative.String())
	}
	return msg.String()

}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) expressionNode() {}
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	var msg bytes.Buffer
	for _, s := range b.Statements {
		msg.WriteString(s.String())
	}
	return msg.String()
}

type FunctionExpression struct {
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionExpression) expressionNode() {}
func (f *FunctionExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionExpression) String() string {
	var msg bytes.Buffer
	param := []string{}
	for _, p := range f.Parameters {
		param = append(param, p.String())
	}

	msg.WriteString(f.TokenLiteral())
	msg.WriteString("(")
	msg.WriteString(strings.Join(param, ", "))
	msg.WriteString(") ")
	msg.WriteString(f.Body.String())

	return msg.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}
func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

func (c *CallExpression) String() string {
	var msg bytes.Buffer

	args := []string{}

	for _, a := range c.Arguments {
		args = append(args, a.String())
	}

	msg.WriteString(c.Function.String())
	msg.WriteString("(")
	msg.WriteString(strings.Join(args, ", "))
	msg.WriteString(")")

	return msg.String()
}

type StringLiteral struct {
	Token   token.Token
	Literal string
}

func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) TokenLiteral() string {
	return s.Token.Literal
}

func (s *StringLiteral) String() string {
	return s.Literal
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode() {}
func (a *ArrayLiteral) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ArrayLiteral) String() string {
	var msg bytes.Buffer
	elementMsg := []string{}

	for _, e := range a.Elements {
		elementMsg = append(elementMsg, e.String())
	}
	msg.WriteString("[")
	msg.WriteString(strings.Join(elementMsg, ", "))
	msg.WriteString("]")

	return msg.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode() {}
func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IndexExpression) String() string {
	var msg bytes.Buffer

	msg.WriteString("(")
	msg.WriteString(i.Left.String())
	msg.WriteString("[")
	msg.WriteString(i.Index.String())
	msg.WriteString("])")

	return msg.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode() {}
func (h *HashLiteral) TokenLiteral() string {
	return h.Token.Literal
}

func (h *HashLiteral) String() string {
	var msg bytes.Buffer

	pairs := []string{}

	for key, val := range h.Pairs {
		pairs = append(pairs, (key.String() + ": " + val.String()))
	}

	msg.WriteString("{")
	msg.WriteString(strings.Join(pairs, ", "))
	msg.WriteString("}")

	return msg.String()
}
