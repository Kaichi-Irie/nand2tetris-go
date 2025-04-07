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
	t   TokenType
	val string
}

var (
	// keywords
	CLASS       Token = Token{t: TT_KEYWORD, val: "class"}
	METHOD      Token = Token{t: TT_KEYWORD, val: "method"}
	FUNCTION    Token = Token{t: TT_KEYWORD, val: "function"}
	CONSTRUCTOR Token = Token{t: TT_KEYWORD, val: "constructor"}
	INT         Token = Token{t: TT_KEYWORD, val: "int"}
	BOOLEAN     Token = Token{t: TT_KEYWORD, val: "boolean"}
	CHAR        Token = Token{t: TT_KEYWORD, val: "char"}
	VOID        Token = Token{t: TT_KEYWORD, val: "void"}
	VAR         Token = Token{t: TT_KEYWORD, val: "var"}
	STATIC      Token = Token{t: TT_KEYWORD, val: "static"}
	FIELD       Token = Token{t: TT_KEYWORD, val: "field"}
	LET         Token = Token{t: TT_KEYWORD, val: "let"}
	DO          Token = Token{t: TT_KEYWORD, val: "do"}
	IF          Token = Token{t: TT_KEYWORD, val: "if"}
	ELSE        Token = Token{t: TT_KEYWORD, val: "else"}
	WHILE       Token = Token{t: TT_KEYWORD, val: "while"}
	RETURN      Token = Token{t: TT_KEYWORD, val: "return"}
	TRUE        Token = Token{t: TT_KEYWORD, val: "true"}
	FALSE       Token = Token{t: TT_KEYWORD, val: "false"}
	NULL        Token = Token{t: TT_KEYWORD, val: "null"}
	THIS        Token = Token{t: TT_KEYWORD, val: "this"}
	// symbols
	DOT        Token = Token{t: TT_SYMBOL, val: "."}
	COMMA      Token = Token{t: TT_SYMBOL, val: ","}
	SEMICOLON  Token = Token{t: TT_SYMBOL, val: ";"}
	UNDERSCORE Token = Token{t: TT_SYMBOL, val: "_"}
	NOT        Token = Token{t: TT_SYMBOL, val: "~"}
	LPAREN     Token = Token{t: TT_SYMBOL, val: "("}
	RPAREN     Token = Token{t: TT_SYMBOL, val: ")"}
	LBRACE     Token = Token{t: TT_SYMBOL, val: "{"}
	RBRACE     Token = Token{t: TT_SYMBOL, val: "}"}
	LSQUARE    Token = Token{t: TT_SYMBOL, val: "["}
	RSQUARE    Token = Token{t: TT_SYMBOL, val: "]"}
	PLUS       Token = Token{t: TT_SYMBOL, val: "+"}
	MINUS      Token = Token{t: TT_SYMBOL, val: "-"}
	ASTERISK   Token = Token{t: TT_SYMBOL, val: "*"}
	SLASH      Token = Token{t: TT_SYMBOL, val: "/"}
	AND        Token = Token{t: TT_SYMBOL, val: "&"}
	OR         Token = Token{t: TT_SYMBOL, val: "|"}
	LESS       Token = Token{t: TT_SYMBOL, val: "<"}
	GREATER    Token = Token{t: TT_SYMBOL, val: ">"}
	EQUAL      Token = Token{t: TT_SYMBOL, val: "="}
)

// Val returns the value of the token
func (t Token) Val() string {
	return t.val
}

// Type returns the type of the token
func (t Token) Type() TokenType {
	return t.t
}

// Is returns true if the token is of the given type. TokenType: TT_KEYWORD, TT_SYMBOL, TT_IDENTIFIER, TT_INT_CONST, TT_STRING_CONST
func (t Token) Is(tType TokenType) bool {
	return t.t == tType
}

// IsPrimitiveType returns true if the token is a primitive type: int, boolean, char or void
func (t Token) IsPrimitiveType() bool {
	switch t.Val() {
	case INT.Val(), BOOLEAN.Val(), CHAR.Val(), VOID.Val():
		return true
	}
	return false
}

// IsOp returns true if the token is an operator: +, -, *, /, &, |, <, >, =
func (t Token) IsOp() bool {
	switch t.Val() {
	case PLUS.Val(), MINUS.Val(), ASTERISK.Val(), SLASH.Val(),
		AND.Val(), OR.Val(), LESS.Val(), GREATER.Val(), EQUAL.Val():
		return true
	}
	return false
}

// IsUnaryOp returns true if the token is a unary operator: ~, -
func (t Token) IsUnaryOp() bool {
	switch t.Val() {
	case NOT.Val(), MINUS.Val():
		return true
	}
	return false
}

// IsKeywordConst returns true if the token is a keyword constant: true, false, null, this
func (t Token) IsKeywordConst() bool {
	switch t.Val() {
	case TRUE.Val(), FALSE.Val(), NULL.Val(), THIS.Val():
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
