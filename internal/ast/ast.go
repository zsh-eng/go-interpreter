package ast

import (
	"bytes"

	"github.com/zsh-eng/go-interpreter/internal/token"
)

type Node interface {
	TokenLiteral() string // For debugging and testing purposes
	String() string
}

type Statement interface {
	Node
	statementNode() // Dummy method to distinguish between statements and expressions
}

type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST our parser produces
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	// Iterate over the statements and append their string representations to a string
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token // The token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	// let <identifier> = <expression>;
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ") // let
	out.WriteString(ls.Name.String())        // <identifier>
	out.WriteString(" = ")                   // =
	if ls.Value != nil {
		out.WriteString(ls.Value.String()) // <expression>
	}
	out.WriteString(";") // ;

	return out.String()
}

type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {} // Expression node as identifiers are expressions
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string {
	return i.Value
}

type ReturnStatement struct {
	Token       token.Token // The token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	// return <expression>;
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ") // return
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String()) // <expression>
	}
	out.WriteString(";") // ;

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // The first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	// <expression>;
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
