package tokenizer

import (
	"fmt"
	"io"
	"nand2tetris-go/vm/vmtranslator"
	"strconv"
	"strings"
)

type TokenType int

const (
	TT_KEYWORD TokenType = iota + 1
	TT_SYMBOL
	TT_IDENTIFIER
	TT_INT_CONST
	TT_STRING_CONST
)

type KeyWordType int

const (
	KT_CLASS KeyWordType = iota + 1
	KT_METHOD
	KT_FUNCTION
	KT_CONSTRUCTOR
	KT_INT
	KT_BOOLEAN
	KT_CHAR
	KT_VOID
	KT_VAR
	KT_STATIC
	KT_FIELD
	KT_LET
	KT_DO
	KT_IF
	KT_ELSE
	KT_WHILE
	KT_RETURN
	KT_TRUE
	KT_FALSE
	KT_NULL
	KT_THIS
)

var XMLEscapes = map[string]string{
	"<":  "&lt;",
	">":  "&gt;",
	"\"": "&quot;",
	"&":  "&amp;",
}

var symbols = []string{
	"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~",
}
var keywords = []string{
	"class", "method", "function", "constructor", "int", "boolean", "char", "void",
	"var", "static", "field", "let", "do", "if", "else", "while", "return",
	"true", "false", "null", "this",
}

// Tokenizer is a struct that reads a file line by line and skips empty lines and comments. It provides a method to get the current line of the scanner without leading and trailing spaces and comments.
type Tokenizer struct {
	Scanner           vmtranslator.CodeScanner
	CurrentLine       string
	CurrentLineLength int
	CurrentPos        int // currentPos is the position of the next token in the current line
	CurrentToken      string
}

// New creates a new Tokenizer with the given reader. It uses a [CodeScanner] to read the file. commentPrefix is the prefix that indicates a comment. Example: "//"
func New(r io.Reader) *Tokenizer {
	t := &Tokenizer{
		Scanner:           vmtranslator.NewCodeScanner(r, "//"),
		CurrentLine:       "",
		CurrentLineLength: 0,
		CurrentPos:        0,
		CurrentToken:      "",
	}
	return t
}

// CreateTokenizerWithFirstToken creates a new Tokenizer with the given reader and advances to the first token. It returns an error if there are no tokens in the file.
func CreateTokenizerWithFirstToken(r io.Reader) (*Tokenizer, error) {
	t := New(r)
	if !t.Advance() {
		return nil, fmt.Errorf("no tokens found")
	}
	return t, nil
}

// Advance advances the scanner to the next token. It returns true if there is a next token, false otherwise.
func (t *Tokenizer) Advance() bool {
	if t.CurrentPos >= t.CurrentLineLength {
		ok := t.Scanner.Scan()
		if !ok {
			return false
		}
		t.CurrentLine = t.Scanner.Text()
		l := len(t.CurrentLine)
		if l == 0 {
			return false
		}
		t.CurrentLineLength = l
		t.CurrentPos = 0
	}

	pos := t.CurrentPos
	// skip spaces
	if t.CurrentLine[pos] == ' ' {
		t.CurrentPos++
		return t.Advance()
	}

	// check if the next token is a symbol
	for _, s := range symbols {
		if pos+len(s) > t.CurrentLineLength {
			continue
		}
		if t.CurrentLine[pos:pos+len(s)] == s {
			t.CurrentPos += len(s)
			t.CurrentToken = s
			return true
		}
	}

	// check if the next token is a keyword
	for _, kw := range keywords {
		if pos+len(kw) > t.CurrentLineLength {
			continue
		}
		if t.CurrentLine[pos:pos+len(kw)] == kw {
			t.CurrentPos += len(kw)
			t.CurrentToken = kw
			return true
		}
	}

	// check if the next token is an integer constant
	if i, ok := ParseIntConst(t.CurrentLine[pos:]); ok == nil {
		t.CurrentPos += len(strconv.Itoa(i))
		t.CurrentToken = strconv.Itoa(i)
		return true
	}

	// check if the next token is a string constant. A string constant is a string that starts and ends with a double quote, so s contains the double quotes.
	if s, ok := ParseStringConst(t.CurrentLine[pos:]); ok == nil {
		fmt.Println("string constant", s)
		t.CurrentPos += len(s)
		t.CurrentToken = s
		return true
	}

	// check if the next token is an identifier
	if id, ok := ParseIdentifier(t.CurrentLine[pos:]); ok == nil {
		t.CurrentPos += len(id)
		t.CurrentToken = id
		return true
	}
	return false
}

// ProcessKeyWord checks if the current token is a keyword of the given type. If it is, it writes the keyword to the writer and advances to the next token. It returns an error if the current token is not a keyword of the given type.
func (t *Tokenizer) ProcessKeyWord(kwType KeyWordType, w io.Writer) error {
	kw := t.CurrentToken
	if kt, err := GetKeyWordType(kw); err != nil {
		return fmt.Errorf("token is not a keyword")
	} else if kt != kwType {
		return fmt.Errorf("token is not the expected keyword")
	}

	_, err := io.WriteString(w, "<keyword> "+kw+" </keyword>")
	if err != nil {
		return err
	}

	t.Advance()
	return nil
}

// ProcessSymbol checks if the current token is a symbol. If it is, it writes the symbol to the writer and advances to the next token. It returns an error if the current token is not a symbol.
func (t *Tokenizer) ProcessSymbol(symbol string, w io.Writer) error {
	if !IsSymbol(symbol) {
		return fmt.Errorf("given symbol is not a valid symbol")
	}
	if symbol != t.CurrentToken {
		return fmt.Errorf("current token is not the expected symbol")
	}

	// escape the symbol
	if escaped, ok := XMLEscapes[symbol]; ok {
		symbol = escaped
	}
	_, err := io.WriteString(w, "<symbol> "+symbol+" </symbol>")
	if err != nil {
		return err
	}
	t.Advance()
	return nil
}

// ProcessIdentifier checks if the current token is an identifier. If it is, it writes the identifier to the writer and advances to the next token. It returns an error if the current token is not an identifier.
func (t *Tokenizer) ProcessIdentifier(w io.Writer) error {
	id := t.CurrentToken
	if !IsIdentifier(id) {
		return fmt.Errorf("token is not an identifier")
	}

	_, err := io.WriteString(w, "<identifier> "+id+" </identifier>")
	if err != nil {
		return err
	}

	t.Advance()
	return nil
}

// GetTokenType returns the type of the token. It returns an error if the token is not a valid token.
func GetTokenType(token string) (TokenType, error) {
	switch {
	case IsIdentifier(token):
		return TT_IDENTIFIER, nil
	case IsIntConst(token):
		return TT_INT_CONST, nil
	case IsStringConst(token):
		return TT_STRING_CONST, nil
	case IsKeyword(token):
		return TT_KEYWORD, nil
	case IsSymbol(token):
		return TT_SYMBOL, nil
	default:
		return 0, fmt.Errorf("unknown token type")
	}
}

// GetKeyWordType returns the keyword type of the token. It returns an error if the token is not a valid keyword.
func GetKeyWordType(token string) (KeyWordType, error) {
	switch token {
	case "class":
		return KT_CLASS, nil
	case "method":
		return KT_METHOD, nil
	case "function":
		return KT_FUNCTION, nil
	case "constructor":
		return KT_CONSTRUCTOR, nil
	case "int":
		return KT_INT, nil
	case "boolean":
		return KT_BOOLEAN, nil
	case "char":
		return KT_CHAR, nil
	case "void":
		return KT_VOID, nil
	case "var":
		return KT_VAR, nil
	case "static":
		return KT_STATIC, nil
	case "field":
		return KT_FIELD, nil
	case "let":
		return KT_LET, nil
	case "do":
		return KT_DO, nil
	case "if":
		return KT_IF, nil
	case "else":
		return KT_ELSE, nil
	case "while":
		return KT_WHILE, nil
	case "return":
		return KT_RETURN, nil
	case "true":
		return KT_TRUE, nil
	case "false":
		return KT_FALSE, nil
	case "null":
		return KT_NULL, nil
	case "this":
		return KT_THIS, nil
	default:
		return 0, fmt.Errorf("unknown keyword type")
	}
}

// func intVal(token string) (int, error)
// func stringVal(token string) (string, error)

// IsSymbol checks if the token is a symbol
func IsSymbol(token string) bool {
	for _, s := range symbols {
		if token == s {
			return true
		}
	}
	return false
}

// IsKeyword checks if the token is a keyword
func IsKeyword(token string) bool {
	for _, kw := range keywords {
		if token == kw {
			return true
		}
	}
	return false
}

func GetIntConst(token string) (int, error) {
	i, err := strconv.Atoi(token)
	if err != nil {
		return 0, fmt.Errorf("not an integer constant")
	}
	return i, nil
}

/*
IsIntConst checks if the string is an integer constant. This does not process any spaces, which means "123 " or "3\n" is not an integer constant
and "123" is an integer constant.
*/
func IsIntConst(token string) bool {
	_, err := GetIntConst(token)
	return err == nil
}

/*
ParseIntConst parses the integer constant at the beginning of the string.
It returns the integer constant and an error if the string does not start with an integer constant.
It also checks if the integer constant is a valid integer constant.
*/
func ParseIntConst(s string) (int, error) {
	// check if the string is integer constant by itself
	if IsIntConst(s) {
		return strconv.Atoi(s)
	}

	// check if the string starts with an integer constant
	i := strings.IndexAny(s, " \t\n\r(){}[];,.+-*/&|<>=~_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if i == -1 {
		i = len(s)
	}
	return GetIntConst(s[0:i])
}

// GetStringConst checks if the string is a string constant and returns the string constant (including the quotes).
func GetStringConst(token string) (string, error) {
	if len(token) < 2 || token[0] != '"' || token[len(token)-1] != '"' {
		return "", fmt.Errorf("not a string constant")
	}
	// check if the string contains any quotes or newlines inside
	if idx := strings.IndexAny(token[1:len(token)-1], "\"\n\r"); idx != -1 {
		return "", fmt.Errorf("string constant contains quotes or newlines")
	}
	return token[0:len(token)], nil
}

// IsStringConst checks if the string is a string constant
func IsStringConst(token string) bool {
	_, err := GetStringConst(token)
	return err == nil
}

/*
ParseStringConst parses the string constant at the beginning of the string.
It returns the string constant (including the quotes) and an error if the string does not start with a string constant.
*/
func ParseStringConst(s string) (string, error) {
	if len(s) < 2 || s[0] != '"' {
		return "", fmt.Errorf("cannot parse string constant")
	}
	idx := strings.Index(s[1:], "\"")
	if idx == -1 {
		return "", fmt.Errorf("cannot find closing quote")
	}

	strConst, err := GetStringConst(s[0 : idx+2])
	if err != nil {
		return "", err
	}
	return strConst, nil
}

// GetIdentifier checks if the string is an identifier. An identifier is a string that starts with a letter or an underscore
func GetIdentifier(token string) (string, error) {
	if len(token) == 0 {
		return "", fmt.Errorf("not an identifier")
	}
	// check if the name does not conflict with a keyword
	if IsKeyword(token) {
		return "", fmt.Errorf("this is a keyword")
	}

	// check if the rest of the string is a letter, digit or underscore
	for i, c := range token {
		// check if the character is a letter, digit or underscore
		if 'A' <= c && c <= 'Z' {
			continue
		}
		if 'a' <= c && c <= 'z' {
			continue
		}
		if i != 0 && '0' <= c && c <= '9' {
			// the first character cannot be a digit
			// but the rest of the string can be a digit
			continue
		}
		if c == '_' {
			continue
		}
		return "", fmt.Errorf("not an identifier. invalid character %c", c)
	}
	return token, nil
}

// IsIdentifier checks if the string is an identifier
func IsIdentifier(token string) bool {
	_, err := GetIdentifier(token)
	return err == nil
}

// ParseIdentifier parses the identifier from the string. An identifier is a string that starts with a letter or an underscore
func ParseIdentifier(s string) (string, error) {
	if IsIdentifier(s) {
		return s, nil
	}
	idx := strings.IndexAny(s, " \t\n\r(){}[];,.+-*/&|<>=~")
	if idx == -1 {
		idx = len(s)
	}
	return GetIdentifier(s[0:idx])
}
