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
	if i, ok := extractIntConst(t.currentLine[pos:]); ok == nil {
		t.currentPos += len(strconv.Itoa(i))
		t.currentToken = strconv.Itoa(i)
		return true
	}

	// check if the next token is a string constant
	if s, ok := extractStringConst(t.currentLine[pos:]); ok == nil {
		t.currentPos += len(s) + 2 // 2 for the quotes
		t.currentToken = s
		return true
	}

	// check if the next token is an identifier
	// TODO: implement identifier extraction
	return false
}

// func getTokenType(token string) tokenType
// func getKeyWordType(token string) keyWordType

// func intVal(token string) (int, error)
// func stringVal(token string) (string, error)

// extractIntConst extracts the integer constant from the string
func extractIntConst(s string) (int, error) {
	i := strings.IndexAny(s, " \t\n\r(){}[];,.+-*/&|<>=~_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if i == -1 || i == 0 {
		return strconv.Atoi(s)
	}
	return strconv.Atoi(s[0:i])
}

// extractStringConst extracts the string constant from the string
func extractStringConst(s string) (string, error) {
	if s[0] != '"' {
		return "", fmt.Errorf("not a string constant")
	}
	idx := strings.Index(s[1:], "\"")
	if idx == -1 {
		return "", fmt.Errorf("not a string constant")
	}
	return s[1 : 1+idx], nil
}
