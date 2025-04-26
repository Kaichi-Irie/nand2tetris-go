package compilationengine

import (
	tk "nand2tetris-go/jackcompiler/tokenizer"
	"strings"
	"testing"
)

func TestProcessKeyWord(t *testing.T) {
	tests := []struct {
		current tk.Token
		kwToken tk.Token
		want    string
	}{

		{tk.Token{
			T:   tk.TT_KEYWORD,
			Val: "class",
		}, tk.CLASS, "<keyword> class </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "method"}, tk.METHOD, "<keyword> method </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "function"}, tk.FUNCTION, "<keyword> function </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "constructor"}, tk.CONSTRUCTOR, "<keyword> constructor </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "int"}, tk.INT, "<keyword> int </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "boolean"}, tk.BOOLEAN, "<keyword> boolean </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "char"}, tk.CHAR, "<keyword> char </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "void"}, tk.VOID, "<keyword> void </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "var"}, tk.VAR, "<keyword> var </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "static"}, tk.STATIC, "<keyword> static </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "field"}, tk.FIELD, "<keyword> field </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "let"}, tk.LET, "<keyword> let </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "do"}, tk.DO, "<keyword> do </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "if"}, tk.IF, "<keyword> if </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "else"}, tk.ELSE, "<keyword> else </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "while"}, tk.WHILE, "<keyword> while </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "return"}, tk.RETURN, "<keyword> return </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "true"}, tk.TRUE, "<keyword> true </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "false"}, tk.FALSE, "<keyword> false </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "null"}, tk.NULL, "<keyword> null </keyword>\n"},
		{tk.Token{T: tk.TT_KEYWORD, Val: "this"}, tk.THIS, "<keyword> this </keyword>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.current.Val, func(t *testing.T) {
			tknz, err := tk.NewWithFirstToken(strings.NewReader(tt.current.Val))
			if err != nil {
				t.Errorf("createTokenizerWithFirstToken(%s) = %v", tt.current.Val, err)
				return
			}
			w := strings.Builder{}
			ce := CompilationEngine{
				t:      tknz,
				writer: &w,
			}
			if err := ce.ProcessKeyWord(tt.kwToken); err != nil {
				t.Errorf("processKeyWord(%s) = %v", tt.current.Val, err)
				return
			}
			if w.String() != tt.want {
				t.Errorf("processKeyWord(%s) = %s, want %s", tt.current.Val, w.String(), tt.want)
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
			w := strings.Builder{}
			ce := *NewWithFirstToken(&w, strings.NewReader(tt.s), "")
			if err := ce.ProcessIdentifier(); err != nil {
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
		symbol tk.Token
		want   string
	}{
		{tk.LBRACE, "<symbol> { </symbol>\n"},
		{tk.RBRACE, "<symbol> } </symbol>\n"},
		{tk.LPAREN, "<symbol> ( </symbol>\n"},
		{tk.RPAREN, "<symbol> ) </symbol>\n"},
		{tk.LSQUARE, "<symbol> [ </symbol>\n"},
		{tk.RSQUARE, "<symbol> ] </symbol>\n"},
		{tk.SEMICOLON, "<symbol> ; </symbol>\n"},
		{tk.COMMA, "<symbol> , </symbol>\n"},
		{tk.DOT, "<symbol> . </symbol>\n"},
		{tk.PLUS, "<symbol> + </symbol>\n"},
		{tk.MINUS, "<symbol> - </symbol>\n"},
		{tk.ASTERISK, "<symbol> * </symbol>\n"},
		{tk.SLASH, "<symbol> / </symbol>\n"},
		{tk.OR, "<symbol> | </symbol>\n"},
		{tk.LESS, "<symbol> &lt; </symbol>\n"},
		{tk.GREATER, "<symbol> &gt; </symbol>\n"},
		{tk.EQUAL, "<symbol> = </symbol>\n"},
		{tk.NOT, "<symbol> ~ </symbol>\n"},
		{tk.AND, "<symbol> &amp; </symbol>\n"},
	}

	for _, tt := range tests {
		t.Run(tt.symbol.Val, func(t *testing.T) {
			tknz, err := tk.NewWithFirstToken(strings.NewReader(tt.symbol.Val))
			if err != nil {
				t.Errorf("createTokenizerWithFirstToken(%s) = %v", tt.symbol.Val, err)
				return
			}
			w := strings.Builder{}
			ce := CompilationEngine{
				t:      tknz,
				writer: &w,
			}
			if err := ce.ProcessSymbol(tt.symbol); err != nil {
				t.Errorf("processSymbol(%s) = %v", tt.symbol.Val, err)
				return
			}
			if w.String() != tt.want {
				t.Errorf("processSymbol(%s) = %s, want %s", tt.symbol.Val, w.String(), tt.want)
			}
		})
	}
}
