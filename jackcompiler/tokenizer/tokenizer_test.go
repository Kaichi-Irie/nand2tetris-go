package tokenizer

import (
	"strings"
	"testing"
)

func TestGetTokenType(t *testing.T) {
	tests := []struct {
		token string
		want  tokenType
	}{
		{"class", TT_KEYWORD},
		{"method", TT_KEYWORD},
		{"function", TT_KEYWORD},
		{"constructor", TT_KEYWORD},
		{"int", TT_KEYWORD},
		{"boolean", TT_KEYWORD},
		{"char", TT_KEYWORD},
		{"void", TT_KEYWORD},
		{"var", TT_KEYWORD},
		{"static", TT_KEYWORD},
		{"field", TT_KEYWORD},
		{"let", TT_KEYWORD},
		{"do", TT_KEYWORD},
		{"if", TT_KEYWORD},
		{"else", TT_KEYWORD},
		{"while", TT_KEYWORD},
		{"return", TT_KEYWORD},
		{"true", TT_KEYWORD},
		{"false", TT_KEYWORD},
		{"null", TT_KEYWORD},
		{"this", TT_KEYWORD},
		{"123", TT_INT_CONST},
		{"0", TT_INT_CONST},
		{"999", TT_INT_CONST},
		{"098", TT_INT_CONST},
		{"hello0", TT_IDENTIFIER},
		{"hello", TT_IDENTIFIER},
		{"_hello", TT_IDENTIFIER},
		{"hello_world", TT_IDENTIFIER},
		{"\"hello\"", TT_STRING_CONST},
		{"\"\"", TT_STRING_CONST},
		{"\"123\"", TT_STRING_CONST},
		{"\"hello world\"", TT_STRING_CONST},
		{"(", TT_SYMBOL},
		{")", TT_SYMBOL},
		{"{", TT_SYMBOL},
		{"}", TT_SYMBOL},
		{"[", TT_SYMBOL},
		{"]", TT_SYMBOL},
		{";", TT_SYMBOL},
		{",", TT_SYMBOL},
		{".", TT_SYMBOL},
		{"+", TT_SYMBOL},
		{"-", TT_SYMBOL},
		{"*", TT_SYMBOL},
		{"/", TT_SYMBOL},
		{"&", TT_SYMBOL},
		{"|", TT_SYMBOL},
		{"<", TT_SYMBOL},
		{">", TT_SYMBOL},
		{"=", TT_SYMBOL},
		{"~", TT_SYMBOL},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got, err := GetTokenType(tt.token); err != nil {
				t.Errorf("getTokenType(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("getTokenType(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestGetKeyWordType(t *testing.T) {
	tests := []struct {
		token string
		want  keyWordType
	}{
		{"class", KT_CLASS},
		{"method", KT_METHOD},
		{"function", KT_FUNCTION},
		{"constructor", KT_CONSTRUCTOR},
		{"int", KT_INT},
		{"boolean", KT_BOOLEAN},
		{"char", KT_CHAR},
		{"void", KT_VOID},
		{"var", KT_VAR},
		{"static", KT_STATIC},
		{"field", KT_FIELD},
		{"let", KT_LET},
		{"do", KT_DO},
		{"if", KT_IF},
		{"else", KT_ELSE},
		{"while", KT_WHILE},
		{"return", KT_RETURN},
		{"true", KT_TRUE},
		{"false", KT_FALSE},
		{"null", KT_NULL},
		{"this", KT_THIS},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got, err := GetKeyWordType(tt.token); err != nil {
				t.Errorf("getKeyWordType(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("getKeyWordType(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestIsIntConst(t *testing.T) {
	tests := []struct {
		token string
		want  bool
	}{
		{"123", true},
		{"0", true},
		{"999", true},
		{"098", true},
		{"098hello", false},
		{"3+5", false},
		{"3   ", false},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got := IsIntConst(tt.token); got != tt.want {
				t.Errorf("isIntConst(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}
func TestExtractIntConst(t *testing.T) {
	tests := []struct {
		token string
		want  int
	}{
		{"123", 123},
		{"0", 0},
		{"999", 999},
		{"098", 98},
		{"098hello", 98},
		{"3+5", 3},
		{"3   ", 3},
		{"4*2", 4},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got, err := ParseIntConst(tt.token); err != nil {
				t.Errorf("extractIntConst(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("extractIntConst(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestIsStringConst(t *testing.T) {
	tests := []struct {
		token string
		want  bool
	}{
		{"\"hello\"", true},
		{"\"\"", true},
		{"\"123\"", true},
		{"\"hello world\"", true},
		{"hello", false},
		{"123", false},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got := IsStringConst(tt.token); got != tt.want {
				t.Errorf("isStringConst(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestExtractStringConst(t *testing.T) {
	tests := []struct {
		token string
		want  string
	}{
		{"\"hello\"", "hello"},
		{"\"\"", ""},
		{"\"123\"", "123"},
		{"\"hello world\"", "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got, err := ParseStringConst(tt.token); err != nil {
				t.Errorf("extractStringConst(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("extractStringConst(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestIsIdentifier(t *testing.T) {
	tests := []struct {
		token string
		want  bool
	}{
		{"hello", true},
		{"_hello", true},
		{"hello_world", true},
		{"hello123", true},
		{"3hello", false},
		{" hello", false},
		{"hello ", false},
		{"hello+world", false},
		{"hello-world", false},
		{"hello.world", false},
		{"hello,world", false},
		{"hello;world", false},
		{"hello:world", false},
		{"hello?world", false},
		{"hello!world", false},
		{"hello@world", false},
		{"hello#world", false},
		{"hello$world", false},
		{"hello%world", false},
		{"hello^world", false},
		{"hello&world", false},
		{"hello*world", false},
		{"hello(world", false},
		{"hello)world", false},
		{"class", false},
		{"method", false},
		{"function", false},
		{"constructor", false},
		{"int", false},
		{"boolean", false},
		{"char", false},
		{"void", false},
		{"var", false},
		{"static", false},
		{"field", false},
		{"let", false},
		{"do", false},
		{"if", false},
		{"else", false},
		{"while", false},
		{"return", false},
		{"true", false},
		{"false", false},
		{"null", false},
		{"this", false},
		{"123", false},
		{"\"hello\"", false},
		{"\"hello world\"", false},
		{"(", false}}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got := IsIdentifier(tt.token); got != tt.want {
				t.Errorf("isIdentifier(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestExtractIdentifier(t *testing.T) {
	tests := []struct {
		token string
		want  string
	}{
		{"hello + 3", "hello"},
		{"_hello(\"str\")", "_hello"},
		{"hello_world", "hello_world"},
		{"hello123", "hello123"},
		{"x)", "x"},
		{"x.y", "x"},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got, err := ParseIdentifier(tt.token); err != nil {
				t.Errorf("extractIdentifier(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("extractIdentifier(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestTokenizer(t *testing.T) {
	var tokenizer = NewTokenizer(strings.NewReader(`
// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/9/Average/Main.jack

/*
comment in multiple lines
*/

/**
this is also supported comment.
*/
class Main {
// comment

// comment
    function void main() {
        var int x;
		var int y;
        let x = Keyboard.readInt("enter the number");
		// comment
        do Output.printInt(x);
	 if (((x+y)<254) & ((x + y)<510)) {
	 			// comment
				do Output.printInt(x);
      } else {
		do Output.printInt(y);
		// comment
}

        	return;

    }// comment
} // comment

`))
	tests := []struct {
		wantToken string
		wantType  tokenType
	}{
		{"class", TT_KEYWORD},
		{"Main", TT_IDENTIFIER},
		{"{", TT_SYMBOL},
		{"function", TT_KEYWORD},
		{"void", TT_KEYWORD},
		{"main", TT_IDENTIFIER},
		{"(", TT_SYMBOL},
		{")", TT_SYMBOL},
		{"{", TT_SYMBOL},
		{"var", TT_KEYWORD},
		{"int", TT_KEYWORD},
		{"x", TT_IDENTIFIER},
		{";", TT_SYMBOL},
		{"var", TT_KEYWORD},
		{"int", TT_KEYWORD},
		{"y", TT_IDENTIFIER},
		{";", TT_SYMBOL},
		{"let", TT_KEYWORD},
		{"x", TT_IDENTIFIER},
		{"=", TT_SYMBOL},
		{"Keyboard.readInt", TT_IDENTIFIER},
		{"(", TT_SYMBOL},
		{"\"enter the number\"", TT_STRING_CONST},
		{")", TT_SYMBOL},
		{";", TT_SYMBOL},
		{"do", TT_KEYWORD},
		{"Output.printInt", TT_IDENTIFIER},
		{"(", TT_SYMBOL},
		{"x", TT_IDENTIFIER},
		{")", TT_SYMBOL},
		{";", TT_SYMBOL},
		{"if", TT_KEYWORD},
		{"(", TT_SYMBOL},
		{"(", TT_SYMBOL},
		{"(", TT_SYMBOL},
		{"x", TT_IDENTIFIER},
		{"+", TT_SYMBOL},
		{"y", TT_IDENTIFIER},
		{")", TT_SYMBOL},
		{"<", TT_SYMBOL},
		{"254", TT_INT_CONST},
		{")", TT_SYMBOL},
		{"&", TT_SYMBOL},
		{"(", TT_SYMBOL},
		{"(", TT_SYMBOL},
		{"x", TT_IDENTIFIER},
		{"+", TT_SYMBOL},
		{"y", TT_IDENTIFIER},
		{")", TT_SYMBOL},
		{"<", TT_SYMBOL},
		{"510", TT_INT_CONST},
		{")", TT_SYMBOL},
		{"{", TT_SYMBOL},
		{"do", TT_KEYWORD},
		{"Output.printInt", TT_IDENTIFIER},
		{"(", TT_SYMBOL},
		{"x", TT_IDENTIFIER},
		{")", TT_SYMBOL},
		{";", TT_SYMBOL},
		{"}", TT_SYMBOL},
		{"else", TT_KEYWORD},
		{"{", TT_SYMBOL},
		{"do", TT_KEYWORD},
		{"Output.printInt", TT_IDENTIFIER},
		{"(", TT_SYMBOL},
		{"y", TT_IDENTIFIER},
		{")", TT_SYMBOL},
		{";", TT_SYMBOL},
		{"}", TT_SYMBOL},
		{"return", TT_KEYWORD},
		{";", TT_SYMBOL},
		{"}", TT_SYMBOL},
		{"}", TT_SYMBOL},
	}

	for i := 0; tokenizer.advance(); i++ {
		token := tokenizer.currentToken
		if i >= len(tests) {
			t.Errorf("tokenizer returned more tokens than expected")
			break
		}
		if tests[i].wantToken != token {
			t.Errorf("tokenizer returned %s, expected %s", token, tests[i].wantToken)
		}
		if tokenType, err := GetTokenType(token); err != nil {
			t.Errorf("getTokenType(%s) returned error: %v", token, err)
		} else if tests[i].wantType != tokenType {
			t.Errorf("tokenizer returned %d, expected %d", tokenType, tests[i].wantType)
		}
	}
}
