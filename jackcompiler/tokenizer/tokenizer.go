package tokenizer

import (
	"fmt"
	"nand2tetris-golang/vm/vmtranslator"
	"strconv"
	"strings"
)

type tokenType int

const (
	TT_KEYWORD tokenType = iota
	TT_SYMBOL
	TT_IDENTIFIER
	TT_INT_CONST
	TT_STRING_CONST
)

type keyWordType int

const (
	KT_CLASS keyWordType = iota
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

/*
IsIntConst checks if the string is an integer constant. This does not process any spaces, which means "123 " or "3\n" is not an integer constant
and "123" is an integer constant.
*/
func IsIntConst(s string) bool {
	_, err := strconv.Atoi(s)
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
	intConst := s[0:i]
	if !IsIntConst(intConst) {
		return 0, fmt.Errorf("not an integer constant")
	}
	return strconv.Atoi(intConst)
}

// IsStringConst checks if the string is a string constant
func IsStringConst(s string) bool {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return false
	}
	// check if the string contains any quotes or newlines inside
	if idx := strings.IndexAny(s[1:len(s)-1], "\"\n\r"); idx != -1 {
		return false
	}
	return true
}

/*
extractStringConst extracts the string constant at the beginning of the string.
It returns the string constant without the quotes and an error if the string is not a string constant.
*/
func ExtractStringConst(s string) (string, error) {
	if len(s) == 0 {
		return "", fmt.Errorf("empty string")
	}
	idx := strings.Index(s[1:], "\"")
	strConst := s[1 : 1+idx]

	if !IsStringConst("\"" + strConst + "\"") {
		return "", fmt.Errorf("not a string constant")
	}
	return strConst, nil
}

func IsIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	// check if the name does not conflict with a keyword
	if IsKeyword(s) {
		return false
	}

	// check if the rest of the string is a letter, digit or underscore
	for i, c := range s {
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
		return false
	}
	return true
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

	identifier := s[0:idx]
	if !IsIdentifier(identifier) {
		return "", fmt.Errorf("not an identifier")
	}
	return identifier, nil
}
