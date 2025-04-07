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
			} else if got.Type() != tt.want.Type() || got.Val() != tt.want.Val() {
				t.Errorf("ParseSymbol(%s) = %v, want %v", tt.s, got.Val(), tt.want.Val())
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
			} else if got.Type() != tt.want.Type() || got.Val() != tt.want.Val() {
				t.Errorf("ParseKeyword(%s) = %v, want %v", tt.s, got.Val(), tt.want.Val())
			}
		})
	}
}

func TestParseStringConst(t *testing.T) {
	tests := []struct {
		s    string
		want Token
	}{
		{"\"hello\"", Token{t: TT_STRING_CONST, val: "\"hello\""}},
		{"\"\"", Token{t: TT_STRING_CONST, val: "\"\""}},
		{"\"123\"", Token{t: TT_STRING_CONST, val: "\"123\""}},
		{"\"hello world\"", Token{t: TT_STRING_CONST, val: "\"hello world\""}},
		{"\"hello\"aaa", Token{t: TT_STRING_CONST, val: "\"hello\""}},
		{"\"hello\"123", Token{t: TT_STRING_CONST, val: "\"hello\""}},
		{"\"hello\"+", Token{t: TT_STRING_CONST, val: "\"hello\""}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got, err := ParseStringConst(tt.s); err != nil {
				t.Errorf("ParseStringConst(%s) = %v, want %v", tt.s, err, tt.want)
			} else if got.Type() != tt.want.Type() || got.Val() != tt.want.Val() {
				t.Errorf("ParseStringConst(%s) = %v, want %v", tt.s, got.Val(), tt.want.Val())
			}
		})
	}
}

func TestParseIdentifier(t *testing.T) {
	tests := []struct {
		s    string
		want Token
	}{
		{"hello + 3", Token{t: TT_IDENTIFIER, val: "hello"}},
		{"hello", Token{t: TT_IDENTIFIER, val: "hello"}},
		{"_hello", Token{t: TT_IDENTIFIER, val: "_hello"}},
		{"_hello(\"str\")", Token{t: TT_IDENTIFIER, val: "_hello"}},
		{"hello_world", Token{t: TT_IDENTIFIER, val: "hello_world"}},
		{"hello123", Token{t: TT_IDENTIFIER, val: "hello123"}},
		{"hello123 + 3", Token{t: TT_IDENTIFIER, val: "hello123"}},
		{"x)", Token{t: TT_IDENTIFIER, val: "x"}},
		{"x.y", Token{t: TT_IDENTIFIER, val: "x"}},
		{"boolx.y.z", Token{t: TT_IDENTIFIER, val: "boolx"}},
		{"returnValue", Token{t: TT_IDENTIFIER, val: "returnValue"}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got, err := ParseIdentifier(tt.s); err != nil {
				t.Errorf("ParseIdentifier(%s) = %v, want %v", tt.s, err, tt.want)
			} else if got.Type() != tt.want.Type() || got.Val() != tt.want.Val() {
				t.Errorf("ParseIdentifier(%s) = %v, want %v", tt.s, got.Val(), tt.want.Val())
			}
		})
	}
}

func TestParseIntConst(t *testing.T) {
	tests := []struct {
		s    string
		want Token
	}{
		{"123 + 3", Token{t: TT_INT_CONST, val: "123"}},
		{"0", Token{t: TT_INT_CONST, val: "0"}},
		{"999", Token{t: TT_INT_CONST, val: "999"}},
		{"098", Token{t: TT_INT_CONST, val: "98"}},
		{"098*3", Token{t: TT_INT_CONST, val: "98"}},
		{"3+5", Token{t: TT_INT_CONST, val: "3"}},
		{"3   ", Token{t: TT_INT_CONST, val: "3"}},
		{"4(x+y)", Token{t: TT_INT_CONST, val: "4"}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got, err := ParseIntConst(tt.s); err != nil {
				t.Errorf("ParseIntConst(%s) = %v, want %v", tt.s, err, tt.want)
			} else if got.Type() != tt.want.Type() || got.Val() != tt.want.Val() {
				t.Errorf("ParseIntConst(%s) = %v, want %v", tt.s, got.Val(), tt.want.Val())
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
		{t: TT_KEYWORD, val: "class"},
		{t: TT_IDENTIFIER, val: "Main"},
		{t: TT_SYMBOL, val: "{"},
		{t: TT_KEYWORD, val: "function"},
		{t: TT_KEYWORD, val: "void"},
		{t: TT_IDENTIFIER, val: "main"},
		{t: TT_SYMBOL, val: "("},
		{t: TT_SYMBOL, val: ")"},
		{t: TT_SYMBOL, val: "{"},
		{t: TT_KEYWORD, val: "var"},
		{t: TT_KEYWORD, val: "int"},
		{t: TT_IDENTIFIER, val: "returnValue"},
		{t: TT_SYMBOL, val: ";"},
		{t: TT_KEYWORD, val: "let"},
		{t: TT_IDENTIFIER, val: "returnValue"},
		{t: TT_SYMBOL, val: "="},
		{t: TT_IDENTIFIER, val: "Keyboard"},
		{t: TT_SYMBOL, val: "."},
		{t: TT_IDENTIFIER, val: "readInt"},
		{t: TT_SYMBOL, val: "("},
		{t: TT_STRING_CONST, val: "\"enter the number\""},
		{t: TT_SYMBOL, val: ")"},
		{t: TT_SYMBOL, val: ";"},
		{t: TT_KEYWORD, val: "do"},
		{t: TT_IDENTIFIER, val: "Output"},
		{t: TT_SYMBOL, val: "."},
		{t: TT_IDENTIFIER, val: "printInt"},
		{t: TT_SYMBOL, val: "("},
		{t: TT_IDENTIFIER, val: "returnValue"},
		{t: TT_SYMBOL, val: ")"},
		{t: TT_SYMBOL, val: ";"},
		{t: TT_KEYWORD, val: "return"},
		{t: TT_SYMBOL, val: ";"},
		{t: TT_SYMBOL, val: "}"},
		{t: TT_SYMBOL, val: "}"},
	}

	for i := 0; tokenizer.Advance(); i++ {
		token := tokenizer.CurrentToken
		if i >= len(tests) {
			t.Errorf("tokenizer returned more tokens than expected")
			break
		}
		if token.Type() != tests[i].Type() {
			t.Errorf("tokenizer returned %d, expected %d", token.Type(), tests[i].Type())
		}
		if token.Val() != tests[i].Val() {
			t.Errorf("tokenizer returned %s, expected %s", token.Val(), tests[i].Val())
		}
	}
}

func TestProcessKeyWord(t *testing.T) {
	tests := []struct {
		current Token
		kwToken Token
		want    string
	}{

		{Token{t: TT_KEYWORD, val: "class"}, CLASS, "<keyword> class </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "method"}, METHOD, "<keyword> method </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "function"}, FUNCTION, "<keyword> function </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "constructor"}, CONSTRUCTOR, "<keyword> constructor </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "int"}, INT, "<keyword> int </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "boolean"}, BOOLEAN, "<keyword> boolean </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "char"}, CHAR, "<keyword> char </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "void"}, VOID, "<keyword> void </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "var"}, VAR, "<keyword> var </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "static"}, STATIC, "<keyword> static </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "field"}, FIELD, "<keyword> field </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "let"}, LET, "<keyword> let </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "do"}, DO, "<keyword> do </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "if"}, IF, "<keyword> if </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "else"}, ELSE, "<keyword> else </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "while"}, WHILE, "<keyword> while </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "return"}, RETURN, "<keyword> return </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "true"}, TRUE, "<keyword> true </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "false"}, FALSE, "<keyword> false </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "null"}, NULL, "<keyword> null </keyword>\n"},
		{Token{t: TT_KEYWORD, val: "this"}, THIS, "<keyword> this </keyword>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.current.Val(), func(t *testing.T) {
			tknz, err := CreateTokenizerWithFirstToken(strings.NewReader(tt.current.Val()))
			if err != nil {
				t.Errorf("createTokenizerWithFirstToken(%s) = %v", tt.current.Val(), err)
				return
			}
			w := strings.Builder{}
			if err := tknz.ProcessKeyWord(tt.kwToken, &w); err != nil {
				t.Errorf("processKeyWord(%s) = %v", tt.current.Val(), err)
				return
			}
			if w.String() != tt.want {
				t.Errorf("processKeyWord(%s) = %s, want %s", tt.current.Val(), w.String(), tt.want)
			}
		})
	}
}

func TestProcessIdentifier(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"hello", "<identifier> hello </identifier>\n"},
		{"_hello", "<identifier> _hello </identifier>\n"},
		{"hello_world", "<identifier> hello_world </identifier>\n"},
		{"hello123", "<identifier> hello123 </identifier>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			tknz, err := CreateTokenizerWithFirstToken(strings.NewReader(tt.s))
			if err != nil {
				t.Errorf("createTokenizerWithFirstToken(%s) = %v", tt.s, err)
				return
			}
			w := strings.Builder{}
			if err := tknz.ProcessIdentifier(&w); err != nil {
				t.Errorf("processIdentifier(%s) = %v", tt.s, err)
				return
			}
			if w.String() != tt.want {
				t.Errorf("processIdentifier(%s) = %s, want %s", tt.s, w.String(), tt.want)
			}
		})
	}
}

func TestProcessSymbol(t *testing.T) {
	tests := []struct {
		symbol Token
		want   string
	}{
		{LBRACE, "<symbol> { </symbol>\n"},
		{RBRACE, "<symbol> } </symbol>\n"},
		{LPAREN, "<symbol> ( </symbol>\n"},
		{RPAREN, "<symbol> ) </symbol>\n"},
		{LSQUARE, "<symbol> [ </symbol>\n"},
		{RSQUARE, "<symbol> ] </symbol>\n"},
		{SEMICOLON, "<symbol> ; </symbol>\n"},
		{COMMA, "<symbol> , </symbol>\n"},
		{DOT, "<symbol> . </symbol>\n"},
		{PLUS, "<symbol> + </symbol>\n"},
		{MINUS, "<symbol> - </symbol>\n"},
		{ASTERISK, "<symbol> * </symbol>\n"},
		{SLASH, "<symbol> / </symbol>\n"},
		{OR, "<symbol> | </symbol>\n"},
		{LESS, "<symbol> &lt; </symbol>\n"},
		{GREATER, "<symbol> &gt; </symbol>\n"},
		{EQUAL, "<symbol> = </symbol>\n"},
		{NOT, "<symbol> ~ </symbol>\n"},
		{AND, "<symbol> &amp; </symbol>\n"},
	}

	for _, tt := range tests {
		t.Run(tt.symbol.Val(), func(t *testing.T) {
			tknz, err := CreateTokenizerWithFirstToken(strings.NewReader(tt.symbol.Val()))
			if err != nil {
				t.Errorf("createTokenizerWithFirstToken(%s) = %v", tt.symbol.Val(), err)
				return
			}
			w := strings.Builder{}
			if err := tknz.ProcessSymbol(tt.symbol, &w); err != nil {
				t.Errorf("processSymbol(%s) = %v", tt.symbol.Val(), err)
				return
			}
			if w.String() != tt.want {
				t.Errorf("processSymbol(%s) = %s, want %s", tt.symbol.Val(), w.String(), tt.want)
			}
		})
	}
}
