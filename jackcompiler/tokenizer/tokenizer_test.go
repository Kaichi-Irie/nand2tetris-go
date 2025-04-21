package tokenizer

import (
	"strings"
	"testing"
)

func TestParseSymbol(t *testing.T) {
	tests := []struct {
		s    string
		want Token
	}{
		{"{8", LBRACE},
		{"};", RBRACE},
		{"(func(", LPAREN},
		{")*5", RPAREN},
		{"[3", LSQUARE},
		{"];", RSQUARE},
		{";\n", SEMICOLON},
		{",\n", COMMA},
		{".hello()", DOT},
		{"+3", PLUS},
		{"-3", MINUS},
		{"*func()", ASTERISK},
		{"/func()", SLASH},
		{"|3", OR},
		{"&4", AND},
		{"<5", LESS},
		{">6", GREATER},
		{"=7", EQUAL},
		{"~8", NOT},
		{"&9", AND}}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got, err := ParseSymbol(tt.s); err != nil {
				t.Errorf("ParseSymbol(%s) = %v, want %v", tt.s, err, tt.want)
			} else if got.T != tt.want.T || got.Val != tt.want.Val {
				t.Errorf("ParseSymbol(%s) = %v, want %v", tt.s, got.Val, tt.want.Val)
			}
		})
	}
}

func TestParseKeyword(t *testing.T) {
	tests := []struct {
		s    string
		want Token
	}{
		{"class Main", CLASS},
		{"method returnValue", METHOD},
		{"function void", FUNCTION},
		{"constructor", CONSTRUCTOR},
		{"int i", INT},
		{"boolean b", BOOLEAN},
		{"char c", CHAR},
		{"void", VOID},
		{"var", VAR},
		{"static", STATIC},
		{"field", FIELD},
		{"let", LET},
		{"do Output()", DO},
		{"if(){}", IF},
		{"else{}", ELSE},
		{"while()", WHILE},
		{"return;", RETURN},
		{"true|", TRUE},
		{"false|", FALSE},
		{"null", NULL},
		{"this", THIS}}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got, err := ParseKeyword(tt.s); err != nil {
				t.Errorf("ParseKeyword(%s) = %v, want %v", tt.s, err, tt.want)
			} else if got.T != tt.want.T || got.Val != tt.want.Val {
				t.Errorf("ParseKeyword(%s) = %v, want %v", tt.s, got.Val, tt.want.Val)
			}
		})
	}
}

func TestParseStringConst(t *testing.T) {
	tests := []struct {
		s    string
		want Token
	}{
		{"\"hello\"", Token{T: TT_STRING_CONST, Val: "\"hello\""}},
		{"\"\"", Token{T: TT_STRING_CONST, Val: "\"\""}},
		{"\"123\"", Token{T: TT_STRING_CONST, Val: "\"123\""}},
		{"\"hello world\"", Token{T: TT_STRING_CONST, Val: "\"hello world\""}},
		{"\"hello\"aaa", Token{T: TT_STRING_CONST, Val: "\"hello\""}},
		{"\"hello\"123", Token{T: TT_STRING_CONST, Val: "\"hello\""}},
		{"\"hello\"+", Token{T: TT_STRING_CONST, Val: "\"hello\""}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got, err := ParseStringConst(tt.s); err != nil {
				t.Errorf("ParseStringConst(%s) = %v, want %v", tt.s, err, tt.want)
			} else if got.T != tt.want.T || got.Val != tt.want.Val {
				t.Errorf("ParseStringConst(%s) = %v, want %v", tt.s, got.Val, tt.want.Val)
			}
		})
	}
}

func TestParseIdentifier(t *testing.T) {
	tests := []struct {
		s    string
		want Token
	}{
		{"hello + 3", Token{T: TT_IDENTIFIER, Val: "hello"}},
		{"hello", Token{T: TT_IDENTIFIER, Val: "hello"}},
		{"_hello", Token{T: TT_IDENTIFIER, Val: "_hello"}},
		{"_hello(\"str\")", Token{T: TT_IDENTIFIER, Val: "_hello"}},
		{"hello_world", Token{T: TT_IDENTIFIER, Val: "hello_world"}},
		{"hello123", Token{T: TT_IDENTIFIER, Val: "hello123"}},
		{"hello123 + 3", Token{T: TT_IDENTIFIER, Val: "hello123"}},
		{"x)", Token{T: TT_IDENTIFIER, Val: "x"}},
		{"x.y", Token{T: TT_IDENTIFIER, Val: "x"}},
		{"boolx.y.z", Token{T: TT_IDENTIFIER, Val: "boolx"}},
		{"returnValue", Token{T: TT_IDENTIFIER, Val: "returnValue"}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got, err := ParseIdentifier(tt.s); err != nil {
				t.Errorf("ParseIdentifier(%s) = %v, want %v", tt.s, err, tt.want)
			} else if got.T != tt.want.T || got.Val != tt.want.Val {
				t.Errorf("ParseIdentifier(%s) = %v, want %v", tt.s, got.Val, tt.want.Val)
			}
		})
	}
}

func TestParseIntConst(t *testing.T) {
	tests := []struct {
		s    string
		want Token
	}{
		{"123 + 3", Token{T: TT_INT_CONST, Val: "123"}},
		{"0", Token{T: TT_INT_CONST, Val: "0"}},
		{"999", Token{T: TT_INT_CONST, Val: "999"}},
		{"098", Token{T: TT_INT_CONST, Val: "98"}},
		{"098*3", Token{T: TT_INT_CONST, Val: "98"}},
		{"3+5", Token{T: TT_INT_CONST, Val: "3"}},
		{"3   ", Token{T: TT_INT_CONST, Val: "3"}},
		{"4(x+y)", Token{T: TT_INT_CONST, Val: "4"}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got, err := ParseIntConst(tt.s); err != nil {
				t.Errorf("ParseIntConst(%s) = %v, want %v", tt.s, err, tt.want)
			} else if got.T != tt.want.T || got.Val != tt.want.Val {
				t.Errorf("ParseIntConst(%s) = %v, want %v", tt.s, got.Val, tt.want.Val)
			}
		})
	}
}

func TestTokenizer(t *testing.T) {
	var tokenizer = New(strings.NewReader(`class Main {
    function void main() {
        var int returnValue;
        let returnValue = Keyboard.readInt("enter the number");
        do Output.printInt(returnValue);
        return;
    }
}`))
	tests := []*Token{
		{T: TT_KEYWORD, Val: "class"},
		{T: TT_IDENTIFIER, Val: "Main"},
		{T: TT_SYMBOL, Val: "{"},
		{T: TT_KEYWORD, Val: "function"},
		{T: TT_KEYWORD, Val: "void"},
		{T: TT_IDENTIFIER, Val: "main"},
		{T: TT_SYMBOL, Val: "("},
		{T: TT_SYMBOL, Val: ")"},
		{T: TT_SYMBOL, Val: "{"},
		{T: TT_KEYWORD, Val: "var"},
		{T: TT_KEYWORD, Val: "int"},
		{T: TT_IDENTIFIER, Val: "returnValue"},
		{T: TT_SYMBOL, Val: ";"},
		{T: TT_KEYWORD, Val: "let"},
		{T: TT_IDENTIFIER, Val: "returnValue"},
		{T: TT_SYMBOL, Val: "="},
		{T: TT_IDENTIFIER, Val: "Keyboard"},
		{T: TT_SYMBOL, Val: "."},
		{T: TT_IDENTIFIER, Val: "readInt"},
		{T: TT_SYMBOL, Val: "("},
		{T: TT_STRING_CONST, Val: "\"enter the number\""},
		{T: TT_SYMBOL, Val: ")"},
		{T: TT_SYMBOL, Val: ";"},
		{T: TT_KEYWORD, Val: "do"},
		{T: TT_IDENTIFIER, Val: "Output"},
		{T: TT_SYMBOL, Val: "."},
		{T: TT_IDENTIFIER, Val: "printInt"},
		{T: TT_SYMBOL, Val: "("},
		{T: TT_IDENTIFIER, Val: "returnValue"},
		{T: TT_SYMBOL, Val: ")"},
		{T: TT_SYMBOL, Val: ";"},
		{T: TT_KEYWORD, Val: "return"},
		{T: TT_SYMBOL, Val: ";"},
		{T: TT_SYMBOL, Val: "}"},
		{T: TT_SYMBOL, Val: "}"},
	}

	for i := 0; tokenizer.Advance(); i++ {
		token := tokenizer.CurrentToken
		if i >= len(tests) {
			t.Errorf("tokenizer returned more tokens than expected")
			break
		}
		if token.T != tests[i].T {
			t.Errorf("tokenizer returned %d, expected %d", token.T, tests[i].T)
		}
		if token.Val != tests[i].Val {
			t.Errorf("tokenizer returned %s, expected %s", token.Val, tests[i].Val)
		}
	}
}
