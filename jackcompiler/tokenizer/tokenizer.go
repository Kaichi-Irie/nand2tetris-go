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
// TODO: implement the tokenizer
	return false
}
