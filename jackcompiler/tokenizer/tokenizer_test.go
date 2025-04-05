package tokenizer

import (
	"strings"
	"testing"
)

func TestGetTokenType(t *testing.T) {
	tests := []struct {
		token string
		want  TokenType
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
		want  KeyWordType
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
		{"\"hello\"", "\"hello\""},
		{"\"\"", "\"\""},
		{"\"123\"aaa", "\"123\""},
		{"\"hello world\"", "\"hello world\""},
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
	var tokenizer = New(strings.NewReader(`class Main {
    function void main() {
        var int x;
        let x = Keyboard.readInt("enter the number");
        do Output.printInt(x);
        return;
    }
}`))
	tests := []struct {
		wantToken string
		wantType  TokenType
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
		{"let", TT_KEYWORD},
		{"x", TT_IDENTIFIER},
		{"=", TT_SYMBOL},
		{"Keyboard", TT_IDENTIFIER},
		{".", TT_SYMBOL},
		{"readInt", TT_IDENTIFIER},
		{"(", TT_SYMBOL},
		{"\"enter the number\"", TT_STRING_CONST},
		{")", TT_SYMBOL},
		{";", TT_SYMBOL},
		{"do", TT_KEYWORD},
		{"Output", TT_IDENTIFIER},
		{".", TT_SYMBOL},
		{"printInt", TT_IDENTIFIER},
		{"(", TT_SYMBOL},
		{"x", TT_IDENTIFIER},
		{")", TT_SYMBOL},
		{";", TT_SYMBOL},
		{"return", TT_KEYWORD},
		{";", TT_SYMBOL},
		{"}", TT_SYMBOL},
		{"}", TT_SYMBOL},
	}

	for i := 0; tokenizer.Advance(); i++ {
		token := tokenizer.CurrentToken
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

func TestProcessKeyWord(t *testing.T) {
	tests := []struct {
		s    string
		T    KeyWordType
		want string
	}{
		{"class", KT_CLASS, "<keyword> class </keyword>\n"},
		{"method", KT_METHOD, "<keyword> method </keyword>\n"},
		{"function", KT_FUNCTION, "<keyword> function </keyword>\n"},
		{"constructor", KT_CONSTRUCTOR, "<keyword> constructor </keyword>\n"},
		{"int", KT_INT, "<keyword> int </keyword>\n"},
		{"boolean", KT_BOOLEAN, "<keyword> boolean </keyword>\n"},
		{"char", KT_CHAR, "<keyword> char </keyword>\n"},
		{"void", KT_VOID, "<keyword> void </keyword>\n"},
		{"var", KT_VAR, "<keyword> var </keyword>\n"},
		{"static", KT_STATIC, "<keyword> static </keyword>\n"},
		{"field", KT_FIELD, "<keyword> field </keyword>\n"},
		{"let", KT_LET, "<keyword> let </keyword>\n"},
		{"do", KT_DO, "<keyword> do </keyword>\n"},
		{"if", KT_IF, "<keyword> if </keyword>\n"},
		{"else", KT_ELSE, "<keyword> else </keyword>\n"},
		{"while", KT_WHILE, "<keyword> while </keyword>\n"},
		{"return", KT_RETURN, "<keyword> return </keyword>\n"},
		{"true", KT_TRUE, "<keyword> true </keyword>\n"},
		{"false", KT_FALSE, "<keyword> false </keyword>\n"},
		{"null", KT_NULL, "<keyword> null </keyword>\n"},
		{"this", KT_THIS, "<keyword> this </keyword>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			tknz, err := CreateTokenizerWithFirstToken(strings.NewReader(tt.s))
			if err != nil {
				t.Errorf("createTokenizerWithFirstToken(%s) = %v", tt.s, err)
				return
			}
			w := strings.Builder{}
			if err := tknz.ProcessKeyWord(tt.T, &w); err != nil {
				t.Errorf("processKeyWord(%s) = %v", tt.s, err)
				return
			}
			if w.String() != tt.want {
				t.Errorf("processKeyWord(%s) = %s, want %s", tt.s, w.String(), tt.want)
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
		s    string
		want string
	}{
		{"{", "<symbol> { </symbol>\n"},
		{"}", "<symbol> } </symbol>\n"},
		{"(", "<symbol> ( </symbol>\n"},
		{")", "<symbol> ) </symbol>\n"},
		{"[", "<symbol> [ </symbol>\n"},
		{"]", "<symbol> ] </symbol>\n"},
		{";", "<symbol> ; </symbol>\n"},
		{",", "<symbol> , </symbol>\n"},
		{".", "<symbol> . </symbol>\n"},
		{"+", "<symbol> + </symbol>\n"},
		{"-", "<symbol> - </symbol>\n"},
		{"*", "<symbol> * </symbol>\n"},
		{"/", "<symbol> / </symbol>\n"},
		{"|", "<symbol> | </symbol>\n"},
		{"<", "<symbol> &lt; </symbol>\n"},
		{">", "<symbol> &gt; </symbol>\n"},
		{"=", "<symbol> = </symbol>\n"},
		{"~", "<symbol> ~ </symbol>\n"},
		{"&", "<symbol> &amp; </symbol>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			tknz, err := CreateTokenizerWithFirstToken(strings.NewReader(tt.s))
			if err != nil {
				t.Errorf("createTokenizerWithFirstToken(%s) = %v", tt.s, err)
				return
			}
			w := strings.Builder{}
			if err := tknz.ProcessSymbol(tt.s, &w); err != nil {
				t.Errorf("processSymbol(%s) = %v", tt.s, err)
				return
			}
			if w.String() != tt.want {
				t.Errorf("processSymbol(%s) = %s, want %s", tt.s, w.String(), tt.want)
			}
		})
	}
}
