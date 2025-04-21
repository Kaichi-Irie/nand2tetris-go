package tokenizer

type TokenType int

const (
	TT_KEYWORD TokenType = iota + 1
	TT_SYMBOL
	TT_IDENTIFIER
	TT_INT_CONST
	TT_STRING_CONST
)

type Token struct {
	T   TokenType
	Val string
}

var (
	// keywords
	CLASS       Token = Token{T: TT_KEYWORD, Val: "class"}
	METHOD      Token = Token{T: TT_KEYWORD, Val: "method"}
	FUNCTION    Token = Token{T: TT_KEYWORD, Val: "function"}
	CONSTRUCTOR Token = Token{T: TT_KEYWORD, Val: "constructor"}
	INT         Token = Token{T: TT_KEYWORD, Val: "int"}
	BOOLEAN     Token = Token{T: TT_KEYWORD, Val: "boolean"}
	CHAR        Token = Token{T: TT_KEYWORD, Val: "char"}
	VOID        Token = Token{T: TT_KEYWORD, Val: "void"}
	VAR         Token = Token{T: TT_KEYWORD, Val: "var"}
	STATIC      Token = Token{T: TT_KEYWORD, Val: "static"}
	FIELD       Token = Token{T: TT_KEYWORD, Val: "field"}
	LET         Token = Token{T: TT_KEYWORD, Val: "let"}
	DO          Token = Token{T: TT_KEYWORD, Val: "do"}
	IF          Token = Token{T: TT_KEYWORD, Val: "if"}
	ELSE        Token = Token{T: TT_KEYWORD, Val: "else"}
	WHILE       Token = Token{T: TT_KEYWORD, Val: "while"}
	RETURN      Token = Token{T: TT_KEYWORD, Val: "return"}
	TRUE        Token = Token{T: TT_KEYWORD, Val: "true"}
	FALSE       Token = Token{T: TT_KEYWORD, Val: "false"}
	NULL        Token = Token{T: TT_KEYWORD, Val: "null"}
	THIS        Token = Token{T: TT_KEYWORD, Val: "this"}
	// symbols
	DOT        Token = Token{T: TT_SYMBOL, Val: "."}
	COMMA      Token = Token{T: TT_SYMBOL, Val: ","}
	SEMICOLON  Token = Token{T: TT_SYMBOL, Val: ";"}
	UNDERSCORE Token = Token{T: TT_SYMBOL, Val: "_"}
	NOT        Token = Token{T: TT_SYMBOL, Val: "~"}
	LPAREN     Token = Token{T: TT_SYMBOL, Val: "("}
	RPAREN     Token = Token{T: TT_SYMBOL, Val: ")"}
	LBRACE     Token = Token{T: TT_SYMBOL, Val: "{"}
	RBRACE     Token = Token{T: TT_SYMBOL, Val: "}"}
	LSQUARE    Token = Token{T: TT_SYMBOL, Val: "["}
	RSQUARE    Token = Token{T: TT_SYMBOL, Val: "]"}
	PLUS       Token = Token{T: TT_SYMBOL, Val: "+"}
	MINUS      Token = Token{T: TT_SYMBOL, Val: "-"}
	ASTERISK   Token = Token{T: TT_SYMBOL, Val: "*"}
	SLASH      Token = Token{T: TT_SYMBOL, Val: "/"}
	AND        Token = Token{T: TT_SYMBOL, Val: "&"}
	OR         Token = Token{T: TT_SYMBOL, Val: "|"}
	LESS       Token = Token{T: TT_SYMBOL, Val: "<"}
	GREATER    Token = Token{T: TT_SYMBOL, Val: ">"}
	EQUAL      Token = Token{T: TT_SYMBOL, Val: "="}
)

// Is returns true if the token is of the given type. TokenType: TT_KEYWORD, TT_SYMBOL, TT_IDENTIFIER, TT_INT_CONST, TT_STRING_CONST
func (t Token) Is(tType TokenType) bool {
	return t.T == tType
}

// IsPrimitiveType returns true if the token is a primitive type: int, boolean, char or void
func (t Token) IsPrimitiveType() bool {
	switch t.Val {
	case INT.Val, BOOLEAN.Val, CHAR.Val, VOID.Val:
		return true
	}
	return false
}

// IsOp returns true if the token is an operator: +, -, *, /, &, |, <, >, =
func (t Token) IsOp() bool {
	switch t.Val {
	case PLUS.Val, MINUS.Val, ASTERISK.Val, SLASH.Val,
		AND.Val, OR.Val, LESS.Val, GREATER.Val, EQUAL.Val:
		return true
	}
	return false
}

// IsUnaryOp returns true if the token is a unary operator: ~, -
func (t Token) IsUnaryOp() bool {
	switch t.Val {
	case NOT.Val, MINUS.Val:
		return true
	}
	return false
}

// IsKeywordConst returns true if the token is a keyword constant: true, false, null, this
func (t Token) IsKeywordConst() bool {
	switch t.Val {
	case TRUE.Val, FALSE.Val, NULL.Val, THIS.Val:
		return true
	}
	return false
}

var Symbols = []Token{
	LPAREN,
	RPAREN,
	LBRACE,
	RBRACE,
	LSQUARE,
	RSQUARE,
	COMMA,
	SEMICOLON,
	DOT,
	PLUS,
	MINUS,
	ASTERISK,
	SLASH,
	AND,
	OR,
	LESS,
	GREATER,
	EQUAL,
	NOT}

var Keywords = []Token{
	CLASS,
	METHOD,
	FUNCTION,
	CONSTRUCTOR,
	INT,
	BOOLEAN,
	CHAR,
	VOID,
	VAR,
	STATIC,
	FIELD,
	LET,
	DO,
	IF,
	ELSE,
	WHILE,
	RETURN,
	TRUE,
	FALSE,
	NULL,
	THIS}
