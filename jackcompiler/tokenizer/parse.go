package tokenizer

import (
	"fmt"
	"strconv"
	"strings"
)

/*
ParseIntConst parses the integer constant at the beginning of the string.
It returns the integer constant and an error if the string does not start with an integer constant.
It also checks if the integer constant is a valid integer constant.
*/
func ParseIntConst(s string) (Token, error) {
	// check if the string starts with an integer constant
	i := strings.IndexAny(s, " \t\n\r(){}[];,.+-*/&|<>=~_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if i == -1 {
		i = len(s)
	}
	val, err := strconv.Atoi(s[0:i])
	if err != nil {
		return Token{}, fmt.Errorf("not an integer constant")
	}
	return Token{t: TT_INT_CONST, val: strconv.Itoa(val)}, nil

}

/*
ParseStringConst parses the string constant at the beginning of the string.
It returns the string constant (including the quotes) and an error if the string does not start with a string constant.
*/
func ParseStringConst(s string) (Token, error) {
	if len(s) < 2 || s[0] != '"' {
		return Token{}, fmt.Errorf("not a string constant")
	}
	idx := strings.Index(s[1:], "\"")
	if idx == -1 {
		return Token{}, fmt.Errorf("cannot find closing quote")
	}
	s = s[0 : idx+2]
	fmt.Println("s:", s)
	// check if the string contains any quotes or newlines inside
	if i := strings.IndexAny(s, "\n\r"); i != -1 {
		return Token{}, fmt.Errorf("string constant contains quotes or newlines")
	}
	return Token{t: TT_STRING_CONST, val: s}, nil
}

// ParseIdentifier parses the identifier from the string. An identifier is a string that starts with a letter or an underscore
func ParseIdentifier(s string) (Token, error) {
	idx := strings.IndexAny(s, " \t\n\r(){}[];,.+-*/&|<>=~")
	if idx == -1 {
		idx = len(s)
	}
	s = s[0:idx]
	// check if the string is empty. if it is empty, then it is not an identifier
	if len(s) == 0 {
		return Token{}, fmt.Errorf("not an identifier")
	}
	// check if the name does not conflict with a keyword
	for _, kw := range Keywords {
		if s == kw.Val() {
			return Token{}, fmt.Errorf("this is a keyword")
		}
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
		return Token{}, fmt.Errorf("not an identifier. invalid character %c", c)
	}

	return Token{t: TT_IDENTIFIER, val: s}, nil
}
