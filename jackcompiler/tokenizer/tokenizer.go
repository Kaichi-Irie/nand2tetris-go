package tokenizer

import (
	"fmt"
	"io"
	"nand2tetris-golang/vm/vmtranslator"
	"strconv"
	"strings"
)

type tokenType int

const (
	TT_KEYWORD tokenType = iota + 1
	TT_SYMBOL
	TT_IDENTIFIER
	TT_INT_CONST
	TT_STRING_CONST
)

type keyWordType int

const (
	KT_CLASS keyWordType = iota + 1
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

var symbols = []string{
	"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~",
}
var keywords = []string{
	"class", "method", "function", "constructor", "int", "boolean", "char", "void",
	"var", "static", "field", "let", "do", "if", "else", "while", "return",
	"true", "false", "null", "this",
}

type Tokenizer struct {
	scanner           vmtranslator.CodeScanner
	currentLine       string
	currentLineLength int
	currentPos        int // currentPos is the position of the next token in the current line
	currentToken      string
}

func (t Tokenizer) advance() bool {
	if t.currentPos >= t.currentLineLength {
		t.currentLine = t.scanner.Text()
		l := len(t.currentLine)
		if l == 0 {
			return false
		}
		t.currentLineLength = l
		t.currentPos = 0
	}

	pos := t.currentPos
	// skip spaces
	if t.currentLine[pos] == ' ' {
		t.currentPos++
		return t.advance()
	}

	// check if the next token is a symbol
	for _, s := range symbols {
		if t.currentLine[pos:pos+len(s)] == s {
			t.currentPos += len(s)
			t.currentToken = s
			return true
		}
	}

	// check if the next token is a keyword
	for _, kw := range keywords {
		if t.currentLine[pos:pos+len(kw)] == kw {
			t.currentPos += len(kw)
			t.currentToken = kw
			return true
		}
	}

	// check if the next token is an integer constant
	if i, ok := ExtractIntConst(t.currentLine[pos:]); ok == nil {
		t.currentPos += len(strconv.Itoa(i))
		t.currentToken = strconv.Itoa(i)
		return true
	}

	// check if the next token is a string constant
	if s, ok := ExtractStringConst(t.currentLine[pos:]); ok == nil {
		t.currentPos += len(s) + 2 // 2 for the quotes
		t.currentToken = s
		return true
	}

	// check if the next token is an identifier
	if id, ok := ExtractIdentifier(t.currentLine[pos:]); ok == nil {
		t.currentPos += len(id)
		t.currentToken = id
		return true
	}
	return false
}

// GetTokenType returns the type of the token. It returns an error if the token is not a valid token.
func GetTokenType(token string) (tokenType, error) {
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
func GetKeyWordType(token string) (keyWordType, error) {
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
ExtractIntConst extracts the integer constant at the beginning of the string.
It returns the integer constant and an error if the string does not start with an integer constant.
It also checks if the integer constant is a valid integer constant.
*/
func ExtractIntConst(s string) (int, error) {
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

// GetStringConst checks if the string is a string constant. A string constant is a string that starts and ends with a double quote
func GetStringConst(token string) (string, error) {
	if len(token) < 2 || token[0] != '"' || token[len(token)-1] != '"' {
		return "", fmt.Errorf("not a string constant")
	}
	// check if the string contains any quotes or newlines inside
	if idx := strings.IndexAny(token[1:len(token)-1], "\"\n\r"); idx != -1 {
		return "", fmt.Errorf("string constant contains quotes or newlines")
	}
	return token[1 : len(token)-1], nil
}

// IsStringConst checks if the string is a string constant
func IsStringConst(token string) bool {
	_, err := GetStringConst(token)
	return err == nil
}

/*
extractStringConst extracts the string constant at the beginning of the string.
It returns the string constant without the quotes and an error if the string is not a string constant.
*/
func ExtractStringConst(s string) (string, error) {
	if len(s) < 2 || s[0] != '"' {
		return "", fmt.Errorf("cannot extract string constant")
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

// ExtractIdentifier extracts the identifier from the string. An identifier is a string that starts with a letter or an underscore
func ExtractIdentifier(s string) (string, error) {
	if IsIdentifier(s) {
		return s, nil
	}
	idx := strings.IndexAny(s, " \t\n\r(){}[];,.+-*/&|<>=~")
	if idx == -1 {
		idx = len(s)
	}
	return GetIdentifier(s[0:idx])
}
