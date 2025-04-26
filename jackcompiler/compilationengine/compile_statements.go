package compilationengine

import (
	"io"
	tk "nand2tetris-go/jackcompiler/tokenizer"
	vw "nand2tetris-go/jackcompiler/vmwriter"
)

func (ce *CompilationEngine) CompileStatements() error {
	_, err := io.WriteString(ce.writer, "<statements>\n")
	if err != nil {
		return err
	}

	// process the statements
LOOP:
	for {
		switch token := ce.t.CurrentToken; {
		case token.Val == tk.LET.Val:
			err = ce.CompileLet()
			if err != nil {
				return err
			}
		case token.Val == tk.IF.Val:
			err = ce.CompileIf()
			if err != nil {
				return err
			}
		case token.Val == tk.WHILE.Val:
			err = ce.CompileWhile()
			if err != nil {
				return err
			}
		case token.Val == tk.DO.Val:
			err = ce.CompileDo()
			if err != nil {
				return err
			}
		case token.Val == tk.RETURN.Val:
			err = ce.CompileReturn()
			if err != nil {
				return err
			}
		default:
			break LOOP
		}
	}

	_, err = io.WriteString(ce.writer, "</statements>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileLet() error {
	_, err := io.WriteString(ce.writer, "<letStatement>\n")
	if err != nil {
		return err
	}

	// process the let keyword
	err = ce.ProcessKeyWord(tk.LET)
	if err != nil {
		return err
	}

	// process the var name
	err = ce.ProcessIdentifier()
	if err != nil {
		return err
	}
	for ce.t.CurrentToken.Val == tk.LSQUARE.Val {
		// process the [
		err = ce.ProcessSymbol(tk.LSQUARE)
		if err != nil {
			return err
		}
		// process the expression
		err = ce.CompileExpression(false)
		if err != nil {
			return err
		}
		// process the ]
		err = ce.ProcessSymbol(tk.RSQUARE)
		if err != nil {
			return err
		}
	}

	// process the =
	err = ce.ProcessSymbol(tk.EQUAL)
	if err != nil {
		return err
	}

	// process the expression
	err = ce.CompileExpression(false)
	if err != nil {
		return err
	}

	// process the ;
	err = ce.ProcessSymbol(tk.SEMICOLON)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.writer, "</letStatement>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileIf() error {
	_, err := io.WriteString(ce.writer, "<ifStatement>\n")
	if err != nil {
		return err
	}
	// process the if keyword
	err = ce.ProcessKeyWord(tk.IF)
	if err != nil {
		return err
	}
	// process the (
	err = ce.ProcessSymbol(tk.LPAREN)
	if err != nil {
		return err
	}
	// process the expression
	err = ce.CompileExpression(false)
	if err != nil {
		return err
	}
	// process the )
	err = ce.ProcessSymbol(tk.RPAREN)
	if err != nil {
		return err
	}
	// process the {
	err = ce.ProcessSymbol(tk.LBRACE)
	if err != nil {
		return err
	}
	// process the statements
	err = ce.CompileStatements()
	if err != nil {
		return err
	}
	// process the }
	err = ce.ProcessSymbol(tk.RBRACE)
	if err != nil {
		return err
	}
	// process the else keyword
	if ce.t.CurrentToken.Val == tk.ELSE.Val {
		err = ce.ProcessKeyWord(tk.ELSE)
		if err != nil {
			return err
		}
		// process the {
		err = ce.ProcessSymbol(tk.LBRACE)
		if err != nil {
			return err
		}
		// process the statements
		err = ce.CompileStatements()
		if err != nil {
			return err
		}
		// process the }
		err = ce.ProcessSymbol(tk.RBRACE)
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(ce.writer, "</ifStatement>\n")
	if err != nil {
		return err
	}
	return nil

}

func (ce *CompilationEngine) CompileWhile() error {
	_, err := io.WriteString(ce.writer, "<whileStatement>\n")
	if err != nil {
		return err
	}

	// process the while keyword
	err = ce.ProcessKeyWord(tk.WHILE)
	if err != nil {
		return err
	}

	// process the (
	err = ce.ProcessSymbol(tk.LPAREN)
	if err != nil {
		return err
	}

	// process the expression
	err = ce.CompileExpression(false)
	if err != nil {
		return err
	}

	// process the )
	err = ce.ProcessSymbol(tk.RPAREN)
	if err != nil {
		return err
	}

	// process the {
	err = ce.ProcessSymbol(tk.LBRACE)
	if err != nil {
		return err
	}

	// process the statements
	err = ce.CompileStatements()
	if err != nil {
		return err
	}

	// process the }
	err = ce.ProcessSymbol(tk.RBRACE)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.writer, "</whileStatement>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileDo() error {
	_, err := io.WriteString(ce.writer, "<doStatement>\n")
	if err != nil {
		return err
	}

	// process the do keyword
	err = ce.ProcessKeyWord(tk.DO)
	if err != nil {
		return err
	}

	// process the subroutine call. Skip the <term> and <expression> tags
	err = ce.CompileExpression(true)
	if err != nil {
		return err
	}

	// process the ;
	err = ce.ProcessSymbol(tk.SEMICOLON)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.writer, "</doStatement>\n")
	if err != nil {
		return err
	}
	// remove the return value from the stack
	ce.vmwriter.WritePop(vw.TEMP, 0)
	return nil
}

func (ce *CompilationEngine) CompileReturn() error {
	_, err := io.WriteString(ce.writer, "<returnStatement>\n")
	if err != nil {
		return err
	}

	// process the return keyword
	err = ce.ProcessKeyWord(tk.RETURN)
	if err != nil {
		return err
	}

	// process the expression. Skip if the case is return;
	if ce.t.CurrentToken.Val != tk.SEMICOLON.Val {
		// process the expression
		err = ce.CompileExpression(false)
		if err != nil {
			return err
		}
	}
	// process the ;
	err = ce.ProcessSymbol(tk.SEMICOLON)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.writer, "</returnStatement>\n")
	if err != nil {
		return err
	}

	// write the return command to the VM writer
	err = ce.vmwriter.WriteReturn(ce.subroutineST.IsCurrentVoidFunc())
	if err != nil {
		return err
	}
	return nil
}
