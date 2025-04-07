package tokenizer

import (
	"fmt"
	"io"
)

// TODO: Move ProcessXXX functions to  compilation engine
// ProcessKeyWord checks if the current token is a keyword of the given type. If it is, it writes the keyword to the writer and advances to the next token. It returns an error if the current token is not a keyword of the given type.
func (t *Tokenizer) ProcessKeyWord(kw Token, w io.Writer) error {
	if !t.CurrentToken.Is(TT_KEYWORD) {
		return fmt.Errorf("token is not a keyword")
	} else if t.CurrentToken.Val() != kw.Val() {
		return fmt.Errorf("token is not the expected keyword")
	}

	_, err := io.WriteString(w, "<keyword> "+kw.Val()+" </keyword>\n")
	if err != nil {
		return err
	}

	t.Advance()
	return nil
}

func (t *Tokenizer) ProcessType(w io.Writer) error {
	// process the type: int, char, boolean or className
	var err error
	token := t.CurrentToken
	switch {
	case token.IsPrimitiveType():
		// process the primitive type
		err = t.ProcessKeyWord(token, w)
		if err != nil {
			return err
		}
	case token.Is(TT_IDENTIFIER):
		err = t.ProcessIdentifier(w)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected int, char, boolean or className, got %v", token.Val())
	}
	return nil
}

// ProcessSymbol checks if the current token is a symbol. If it is, it writes the symbol to the writer and advances to the next token. It returns an error if the current token is not a symbol.
func (t *Tokenizer) ProcessSymbol(symbol Token, w io.Writer) error {
	token := t.CurrentToken
	if !token.Is(TT_SYMBOL) {
		return fmt.Errorf("token is not a symbol")
	} else if token.Val() != symbol.Val() {
		return fmt.Errorf("current token is not the expected symbol")
	}

	val := token.Val()
	// escape the symbol
	if escaped, ok := XMLEscapes[val]; ok {
		val = escaped
	}
	_, err := io.WriteString(w, "<symbol> "+val+" </symbol>\n")
	if err != nil {
		return err
	}
	t.Advance()
	return nil
}

// ProcessIdentifier checks if the current token is an identifier. If it is, it writes the identifier to the writer and advances to the next token. It returns an error if the current token is not an identifier.
func (t *Tokenizer) ProcessIdentifier(w io.Writer) error {
	token := t.CurrentToken
	if !token.Is(TT_IDENTIFIER) {
		return fmt.Errorf("token is not an identifier")
	}

	_, err := io.WriteString(w, "<identifier> "+token.Val()+" </identifier>\n")
	if err != nil {
		return err
	}

	t.Advance()
	return nil
}

// ProcessStringConst checks if the current token is a string constant. If it is, it writes the string constant to the writer and advances to the next token. It returns an error if the current token is not a string constant.
func (t *Tokenizer) ProcessStringConst(w io.Writer) error {
	token := t.CurrentToken
	if !token.Is(TT_STRING_CONST) {
		return fmt.Errorf("token is not a string constant")
	}
	// remove the quotes from the string constant
	trimmedStrConst := token.Val()[1 : len(token.Val())-1]
	_, err := io.WriteString(w, "<stringConstant> "+trimmedStrConst+" </stringConstant>\n")
	if err != nil {
		return err
	}

	t.Advance()
	return nil
}

// ProcessIntConst checks if the current token is an integer constant. If it is, it writes the integer constant to the writer and advances to the next token. It returns an error if the current token is not an integer constant.
func (t *Tokenizer) ProcessIntConst(w io.Writer) error {
	token := t.CurrentToken
	if !token.Is(TT_INT_CONST) {
		return fmt.Errorf("token is not an integer constant")
	}

	_, err := io.WriteString(w, "<integerConstant> "+token.Val()+" </integerConstant>\n")
	if err != nil {
		return err
	}

	t.Advance()
	return nil
}
