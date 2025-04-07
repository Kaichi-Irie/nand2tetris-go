package compilationengine

import (
	"io"
	"nand2tetris-go/jackcompiler/tokenizer"
)

func (ce *CompilationEngine) CompileStatements() error {
	_, err := io.WriteString(ce.xmlFile, "<statements>\n")
	if err != nil {
		return err
	}

	// process the statements
LOOP:
	for {
		switch token := ce.t.CurrentToken; {
		case token.Val() == tokenizer.LET.Val():
			err = ce.CompileLet()
			if err != nil {
				return err
			}
		case token.Val() == tokenizer.IF.Val():
			err = ce.CompileIf()
			if err != nil {
				return err
			}
		case token.Val() == tokenizer.WHILE.Val():
			err = ce.CompileWhile()
			if err != nil {
				return err
			}
		case token.Val() == tokenizer.DO.Val():
			err = ce.CompileDo()
			if err != nil {
				return err
			}
		case token.Val() == tokenizer.RETURN.Val():
			err = ce.CompileReturn()
			if err != nil {
				return err
			}
		default:
			break LOOP
		}
	}

	_, err = io.WriteString(ce.xmlFile, "</statements>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileLet() error {
	_, err := io.WriteString(ce.xmlFile, "<letStatement>\n")
	if err != nil {
		return err
	}

	// process the let keyword
	err = ce.t.ProcessKeyWord(tokenizer.LET, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the var name
	err = ce.t.ProcessIdentifier(ce.xmlFile)
	if err != nil {
		return err
	}
	for ce.t.CurrentToken.Val() == tokenizer.LSQUARE.Val() {
		// process the [
		err = ce.t.ProcessSymbol(tokenizer.LSQUARE, ce.xmlFile)
		if err != nil {
			return err
		}
		// process the expression
		err = ce.CompileExpression(false)
		if err != nil {
			return err
		}
		// process the ]
		err = ce.t.ProcessSymbol(tokenizer.RSQUARE, ce.xmlFile)
		if err != nil {
			return err
		}
	}

	// process the =
	err = ce.t.ProcessSymbol(tokenizer.EQUAL, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the expression
	err = ce.CompileExpression(false)
	if err != nil {
		return err
	}

	// process the ;
	err = ce.t.ProcessSymbol(tokenizer.SEMICOLON, ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</letStatement>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileIf() error {
	_, err := io.WriteString(ce.xmlFile, "<ifStatement>\n")
	if err != nil {
		return err
	}
	// process the if keyword
	err = ce.t.ProcessKeyWord(tokenizer.IF, ce.xmlFile)
	if err != nil {
		return err
	}
	// process the (
	err = ce.t.ProcessSymbol(tokenizer.LPAREN, ce.xmlFile)
	if err != nil {
		return err
	}
	// process the expression
	err = ce.CompileExpression(false)
	if err != nil {
		return err
	}
	// process the )
	err = ce.t.ProcessSymbol(tokenizer.RPAREN, ce.xmlFile)
	if err != nil {
		return err
	}
	// process the {
	err = ce.t.ProcessSymbol(tokenizer.LBRACE, ce.xmlFile)
	if err != nil {
		return err
	}
	// process the statements
	err = ce.CompileStatements()
	if err != nil {
		return err
	}
	// process the }
	err = ce.t.ProcessSymbol(tokenizer.RBRACE, ce.xmlFile)
	if err != nil {
		return err
	}
	// process the else keyword
	if ce.t.CurrentToken.Val() == tokenizer.ELSE.Val() {
		err = ce.t.ProcessKeyWord(tokenizer.ELSE, ce.xmlFile)
		if err != nil {
			return err
		}
		// process the {
		err = ce.t.ProcessSymbol(tokenizer.LBRACE, ce.xmlFile)
		if err != nil {
			return err
		}
		// process the statements
		err = ce.CompileStatements()
		if err != nil {
			return err
		}
		// process the }
		err = ce.t.ProcessSymbol(tokenizer.RBRACE, ce.xmlFile)
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(ce.xmlFile, "</ifStatement>\n")
	if err != nil {
		return err
	}
	return nil

}

func (ce *CompilationEngine) CompileWhile() error {
	_, err := io.WriteString(ce.xmlFile, "<whileStatement>\n")
	if err != nil {
		return err
	}

	// process the while keyword
	err = ce.t.ProcessKeyWord(tokenizer.WHILE, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the (
	err = ce.t.ProcessSymbol(tokenizer.LPAREN, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the expression
	err = ce.CompileExpression(false)
	if err != nil {
		return err
	}

	// process the )
	err = ce.t.ProcessSymbol(tokenizer.RPAREN, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the {
	err = ce.t.ProcessSymbol(tokenizer.LBRACE, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the statements
	err = ce.CompileStatements()
	if err != nil {
		return err
	}

	// process the }
	err = ce.t.ProcessSymbol(tokenizer.RBRACE, ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</whileStatement>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileDo() error {
	_, err := io.WriteString(ce.xmlFile, "<doStatement>\n")
	if err != nil {
		return err
	}

	// process the do keyword
	err = ce.t.ProcessKeyWord(tokenizer.DO, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the subroutine call. Skip the <term> and <expression> tags
	err = ce.CompileExpression(true)
	if err != nil {
		return err
	}

	// process the ;
	err = ce.t.ProcessSymbol(tokenizer.SEMICOLON, ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</doStatement>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileReturn() error {
	_, err := io.WriteString(ce.xmlFile, "<returnStatement>\n")
	if err != nil {
		return err
	}

	// process the return keyword
	err = ce.t.ProcessKeyWord(tokenizer.RETURN, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the expression. Skip if the case is return;
	if ce.t.CurrentToken.Val() != tokenizer.SEMICOLON.Val() {
		// process the expression
		err = ce.CompileExpression(false)
		if err != nil {
			return err
		}
	}
	// process the ;
	err = ce.t.ProcessSymbol(tokenizer.SEMICOLON, ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</returnStatement>\n")
	if err != nil {
		return err
	}
	return nil
}
