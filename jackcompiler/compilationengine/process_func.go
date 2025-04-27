package compilationengine

import (
	"fmt"
	tk "nand2tetris-go/jackcompiler/tokenizer"
)

// ProcessKeyWord checks if the current token is a keyword (class, method, function, constructor, int, char, boolean, void, var, static, field, let, do, if, else, while, return, true, false, null, this) and if it matches the expected keyword. If it does, it writes the keyword to the writer and advances to the next token. It returns an error if the current token is not a keyword or does not match the expected keyword.
func (ce *CompilationEngine) ProcessKeyWord(kw tk.Token) error {
	token := ce.t.CurrentToken
	if !token.Is(tk.TT_KEYWORD) {
		return fmt.Errorf("token is not a keyword")
	} else if token.Val != kw.Val {
		return fmt.Errorf("token is not the expected keyword")
	}
	ce.t.Advance()
	return nil
}

// process the type: int, char, boolean or className
func (ce *CompilationEngine) ProcessType() error {
	var err error
	token := ce.t.CurrentToken
	switch {
	case token.IsPrimitiveType():
		// process the primitive type
		err = ce.ProcessKeyWord(token)
		if err != nil {
			return err
		}
	case token.Is(tk.TT_IDENTIFIER):
		err = ce.ProcessIdentifier()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected int, char, boolean or className, got %v", token.Val)
	}
	return nil
}

// ProcessSymbol checks if the current token is a symbol. If it is, it writes the symbol to the writer and advances to the next token. It returns an error if the current token is not a symbol.
func (ce *CompilationEngine) ProcessSymbol(symbol tk.Token) error {
	token := ce.t.CurrentToken
	if !token.Is(tk.TT_SYMBOL) {
		return fmt.Errorf("token is not a symbol")
	} else if token.Val != symbol.Val {
		return fmt.Errorf("current token is not the expected symbol")
	}
	// escape the symbol
	ce.t.Advance()
	return nil
}

// ProcessIdentifier checks if the current token is an identifier. If it is, it writes the identifier to the writer and advances to the next token. It returns an error if the current token is not an identifier.
func (ce *CompilationEngine) ProcessIdentifier() error {
	token := ce.t.CurrentToken
	if !token.Is(tk.TT_IDENTIFIER) {
		return fmt.Errorf("token is not an identifier")
	}
	ce.t.Advance()
	return nil
}

// ProcessStringConst checks if the current token is a string constant. If it is, it writes the string constant to the writer and advances to the next token. It returns an error if the current token is not a string constant.
func (ce *CompilationEngine) ProcessStringConst() error {
	token := ce.t.CurrentToken
	if !token.Is(tk.TT_STRING_CONST) {
		return fmt.Errorf("token is not a string constant")
	}
	ce.t.Advance()
	return nil
}

// ProcessIntConst checks if the current token is an integer constant. If it is, it writes the integer constant to the ce.writer  and advances to the next token. It returns an error if the current token is not an integer constant.
func (ce *CompilationEngine) ProcessIntConst() error {
	token := ce.t.CurrentToken
	if !token.Is(tk.TT_INT_CONST) {
		return fmt.Errorf("token is not an integer constant")
	}
	ce.t.Advance()
	return nil
}
