package compilationengine

import (
	"bytes"
	"strings"
	"testing"

	tk "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/tokenizer"
)

func TestProcessKeyWord(t *testing.T) {
	tests := []struct {
		name        string
		jackCode    string
		expectedKW  tk.Token
		expectError bool
		nextToken   tk.Token // Expected token after successful processing
	}{
		{"Correct keyword", "class Test", tk.CLASS, false, tk.Token{
			T:   tk.TT_IDENTIFIER,
			Val: "Test",
		}},
		{"Incorrect keyword", "method Test", tk.CLASS, true, tk.Token{}},
		{"Not a keyword", "Test class", tk.CLASS, true, tk.Token{}},
		{"Correct keyword (int)", "int i", tk.INT, false, tk.Token{
			T:   tk.TT_IDENTIFIER,
			Val: "i",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithFirstToken(vmFile, strings.NewReader(tt.jackCode), "Test")

			err := ce.ProcessKeyWord(tt.expectedKW)

			if tt.expectError {
				if err == nil {
					t.Errorf("ProcessKeyWord(%v) expected error, got nil", tt.expectedKW)
				}
			} else {
				if err != nil {
					t.Errorf("ProcessKeyWord(%v) unexpected error: %v", tt.expectedKW, err)
				}
				// Check if tokenizer advanced
				if ce.t.CurrentToken.T != tt.nextToken.T {
					t.Errorf("After ProcessKeyWord(%v), expected next token type %v, got %v", tt.expectedKW, tt.nextToken.T, ce.t.CurrentToken.T)
				}
			}
		})
	}
}

func TestProcessType(t *testing.T) {
	tests := []struct {
		name        string
		jackCode    string
		expectError bool
		nextToken   tk.Token // Expected token after successful processing
	}{
		{"Primitive type int", "int varName", false, tk.Token{T: tk.TT_IDENTIFIER, Val: "varName"}},
		{"Primitive type char", "char varName", false, tk.Token{T: tk.TT_IDENTIFIER, Val: "varName"}},
		{"Primitive type boolean", "boolean varName", false, tk.Token{T: tk.TT_IDENTIFIER, Val: "varName"}},
		{"Class name type", "MyClass varName", false, tk.Token{T: tk.TT_IDENTIFIER, Val: "varName"}},
		{"Invalid type (symbol)", "{ varName", true, tk.Token{}},
		{"Invalid type (keyword)", "let varName", true, tk.Token{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithFirstToken(vmFile, strings.NewReader(tt.jackCode), "Test")

			err := ce.ProcessType()

			if tt.expectError {
				if err == nil {
					t.Errorf("ProcessType() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ProcessType() unexpected error: %v", err)
				}
				// Check if tokenizer advanced
				if ce.t.CurrentToken.T != tt.nextToken.T {
					t.Errorf("After ProcessType(), expected next token type %v, got %v", tt.nextToken.T, ce.t.CurrentToken.T)
				}
			}
		})
	}
}

func TestProcessSymbol(t *testing.T) {
	tests := []struct {
		name           string
		jackCode       string
		expectedSymbol tk.Token
		expectError    bool
		nextToken      tk.Token // Expected token after successful processing
	}{
		{"Correct symbol {", "{ var", tk.LBRACE, false, tk.VAR},
		{"Correct symbol ;", "; }", tk.SEMICOLON, false, tk.RBRACE},
		{"Incorrect symbol", "( var", tk.LBRACE, true, tk.Token{}},
		{"Not a symbol", "var {", tk.LBRACE, true, tk.Token{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithFirstToken(vmFile, strings.NewReader(tt.jackCode), "Test")

			err := ce.ProcessSymbol(tt.expectedSymbol)

			if tt.expectError {
				if err == nil {
					t.Errorf("ProcessSymbol(%v) expected error, got nil", tt.expectedSymbol)
				}
			} else {
				if err != nil {
					t.Errorf("ProcessSymbol(%v) unexpected error: %v", tt.expectedSymbol, err)
				}
				// Check if tokenizer advanced
				if ce.t.CurrentToken.T != tt.nextToken.T {
					t.Errorf("After ProcessSymbol(%v), expected next token type %v, got %v", tt.expectedSymbol, tt.nextToken.T, ce.t.CurrentToken.T)
				}
			}
		})
	}
}

func TestProcessIdentifier(t *testing.T) {
	tests := []struct {
		name        string
		jackCode    string
		expectError bool
		nextToken   tk.Token // Expected token after successful processing
	}{
		{"Correct identifier", "myVar ;", false, tk.SEMICOLON},
		{"Not an identifier (keyword)", "class Test", true, tk.Token{}},
		{"Not an identifier (symbol)", "{ Test", true, tk.Token{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithFirstToken(vmFile, strings.NewReader(tt.jackCode), "Test")

			err := ce.ProcessIdentifier()

			if tt.expectError {
				if err == nil {
					t.Errorf("ProcessIdentifier() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ProcessIdentifier() unexpected error: %v", err)
				}
				// Check if tokenizer advanced
				if ce.t.CurrentToken.T != tt.nextToken.T {
					t.Errorf("After ProcessIdentifier(), expected next token type %v, got %v", tt.nextToken.T, ce.t.CurrentToken.T)
				}
			}
		})
	}
}

func TestProcessStringConst(t *testing.T) {
	tests := []struct {
		name        string
		jackCode    string
		expectError bool
		nextToken   tk.Token // Expected token after successful processing
	}{
		{"Correct string const", `"hello" ;`, false, tk.SEMICOLON},
		{"Not a string const (identifier)", `hello ;`, true, tk.Token{}},
		{"Not a string const (int)", `123 ;`, true, tk.Token{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithFirstToken(vmFile, strings.NewReader(tt.jackCode), "Test")

			err := ce.ProcessStringConst()

			if tt.expectError {
				if err == nil {
					t.Errorf("ProcessStringConst() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ProcessStringConst() unexpected error: %v", err)
				}
				// Check if tokenizer advanced
				if ce.t.CurrentToken.T != tt.nextToken.T {
					t.Errorf("After ProcessStringConst(), expected next token type %v, got %v", tt.nextToken.T, ce.t.CurrentToken.T)
				}
			}
		})
	}
}

func TestProcessIntConst(t *testing.T) {
	tests := []struct {
		name        string
		jackCode    string
		expectError bool
		nextToken   tk.Token // Expected token after successful processing
	}{
		{"Correct int const", `123 ;`, false, tk.SEMICOLON},
		{"Not an int const (identifier)", `num ;`, true, tk.Token{}},
		{"Not an int const (string)", `"123" ;`, true, tk.Token{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithFirstToken(vmFile, strings.NewReader(tt.jackCode), "Test")

			err := ce.ProcessIntConst()

			if tt.expectError {
				if err == nil {
					t.Errorf("ProcessIntConst() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ProcessIntConst() unexpected error: %v", err)
				}
				// Check if tokenizer advanced
				if ce.t.CurrentToken.T != tt.nextToken.T {
					t.Errorf("After ProcessIntConst(), expected next token type %v, got %v", tt.nextToken.T, ce.t.CurrentToken.T)
				}
			}
		})
	}
}
