package parser

import (
	"testing"

	"github.com/zsh-eng/go-interpreter/internal/ast"
	"github.com/zsh-eng/go-interpreter/internal/lexer"
	"github.com/zsh-eng/go-interpreter/internal/token"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	// Test the let statements
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	// Iterate over the statements and test each one
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return // Test failed
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement) // Type assertion
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not %s. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	// Test the return statements
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement) // Type assertion

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() not 'return'. got=%q", returnStmt.TokenLiteral())
		}
	}

}

func TestString(t *testing.T) {
	out := `let myVar = anotherVar;`

	// Construct an AST manually and compare its string representation to the expected output
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != out {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return // No errors
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
