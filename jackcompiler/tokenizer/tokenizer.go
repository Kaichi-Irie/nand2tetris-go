package tokenizer

import (
	"fmt"
	"io"
	"nand2tetris-go/vm/vmtranslator"
)

// Tokenizer is a struct that reads a file line by line and skips empty lines and comments. It provides a method to get the current line of the scanner without leading and trailing spaces and comments.
type Tokenizer struct {
	Scanner           vmtranslator.CodeScanner
	CurrentLine       string
	CurrentLineLength int
	CurrentPos        int // currentPos is the position of the next token in the current line
	CurrentToken      Token
}

// New creates a new Tokenizer with the given reader. It uses a [CodeScanner] to read the file. commentPrefix is the prefix that indicates a comment. Example: "//"
func New(r io.Reader) *Tokenizer {
	t := &Tokenizer{
		Scanner:           vmtranslator.NewCodeScanner(r, "//"),
		CurrentLine:       "",
		CurrentLineLength: 0,
		CurrentPos:        0,
		CurrentToken:      Token{},
	}
	return t
}

// CreateTokenizerWithFirstToken creates a new Tokenizer with the given reader and advances to the first token. It returns an error if there are no tokens in the file.
func CreateTokenizerWithFirstToken(r io.Reader) (*Tokenizer, error) {
	t := New(r)
	if !t.Advance() {
		return nil, fmt.Errorf("create tokenizer with first token: no tokens in the file")
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
	for _, s := range Symbols {
		length := len(s.Val())
		if pos+length > t.CurrentLineLength {
			continue
		}
		if t.CurrentLine[pos:pos+length] == s.Val() {
			t.CurrentPos += length
			t.CurrentToken = s
			return true
		}
	}

	// check if the next token is a keyword
	for _, kw := range Keywords {
		length := len(kw.Val())
		if pos+length > t.CurrentLineLength {
			continue
		}

		kwCandidate := t.CurrentLine[pos : pos+length]
		if pos+length < t.CurrentLineLength {
			followingChar := t.CurrentLine[pos+length]
			// if followingChar is alphanumeric, or underscore, then it is not a keyword. followingChar is the character just after the keyword
			if 'a' <= followingChar && followingChar <= 'z' ||
				'A' <= followingChar && followingChar <= 'Z' ||
				'0' <= followingChar && followingChar <= '9' ||
				followingChar == '_' {
				continue
			}
		}
		if kwCandidate == kw.Val() {
			t.CurrentPos += length
			t.CurrentToken = kw
			return true
		}
	}

	// check if the next token is an integer constant
	if i, ok := ParseIntConst(t.CurrentLine[pos:]); ok == nil {
		t.CurrentPos += len(i.Val())
		t.CurrentToken = i
		return true
	}

	// check if the next token is a string constant. A string constant is a string that starts and ends with a double quote, so s contains the double quotes.
	if s, ok := ParseStringConst(t.CurrentLine[pos:]); ok == nil {
		t.CurrentPos += len(s.Val())
		t.CurrentToken = s
		return true
	}

	// check if the next token is an identifier
	if id, ok := ParseIdentifier(t.CurrentLine[pos:]); ok == nil {
		t.CurrentPos += len(id.Val())
		t.CurrentToken = id
		return true
	}
	return false
}
