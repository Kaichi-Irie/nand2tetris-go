package compilationengine

import (
	"fmt"
	tk "nand2tetris-go/jackcompiler/tokenizer"
	vw "nand2tetris-go/jackcompiler/vmwriter"
	"strconv"
)

func (ce *CompilationEngine) CompileStatements() error {
	var err error
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
	return nil
}

func (ce *CompilationEngine) CompileLet() error {
	// process the let keyword
	err := ce.ProcessKeyWord(tk.LET)
	if err != nil {
		return err
	}

	// process the var name
	varName := ce.t.CurrentToken.Val
	err = ce.ProcessIdentifier()
	if err != nil {
		return err
	}
	// process the [
	isArray := false
	if ce.t.CurrentToken.Val == tk.LSQUARE.Val {
		isArray = true
		// push array base address
		id, ok := ce.Lookup(varName)
		if !ok {
			return fmt.Errorf("variable %s is not defined. LetStatement cannot be used", varName)
		}
		seg := vw.SegmentOfKind[id.Kind]
		index := id.Index
		err = ce.vmwriter.WritePush(seg, index)
		if err != nil {
			return err
		}

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

		// add
		err = ce.vmwriter.WriteArithmetic(vw.ADD)
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

	if isArray {
		// pop the destination expression to temp segment
		err = ce.vmwriter.WritePop(vw.TEMP, 0)
		if err != nil {
			return err
		}
		err = ce.vmwriter.WritePop(vw.POINTER, 1)
		if err != nil {
			return err
		}
		err = ce.vmwriter.WritePush(vw.TEMP, 0)
		if err != nil {
			return err
		}
		err = ce.vmwriter.WritePop(vw.THAT, 0)
		if err != nil {
			return err
		}
		// pop the value to the variable
	} else if identifier, ok := ce.Lookup(varName); ok {
		seg := vw.SegmentOfKind[identifier.Kind]
		index := identifier.Index
		err = ce.vmwriter.WritePop(seg, index)
		if err != nil {
			return err
		}
	} else {
		// TODO: handle the case where the variable is not defined
		// return fmt.Errorf("variable %s is not defined. LetStatement cannot be used", varName)
	}
	return nil
}

func (ce *CompilationEngine) CompileIf() error {
	// process the if keyword
	err := ce.ProcessKeyWord(tk.IF)
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

	// not
	err = ce.vmwriter.WriteArithmetic(vw.NOT)
	if err != nil {
		return err
	}

	// if-goto labelElse
	labelElse := "label" + strconv.Itoa(ce.labelCount)
	err = ce.vmwriter.WriteIf(labelElse)
	if err != nil {
		return err
	}
	ce.labelCount++

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

	// goto labelFinally
	labelFinally := "label" + strconv.Itoa(ce.labelCount)
	err = ce.vmwriter.WriteGoto(labelFinally)
	if err != nil {
		return err
	}
	ce.labelCount++

	// process the }
	err = ce.ProcessSymbol(tk.RBRACE)
	if err != nil {
		return err
	}

	// labelElse
	err = ce.vmwriter.WriteLabel(labelElse)
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

	// labelFinally
	err = ce.vmwriter.WriteLabel(labelFinally)
	if err != nil {
		return err
	}
	return nil

}

func (ce *CompilationEngine) CompileWhile() error {
	// labelWhile
	labelWhile := "label" + strconv.Itoa(ce.labelCount)
	err := ce.vmwriter.WriteLabel(labelWhile)
	if err != nil {
		return err
	}
	ce.labelCount++

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

	// not
	err = ce.vmwriter.WriteArithmetic(vw.NOT)
	if err != nil {
		return err
	}
	// if-goto labelFinally
	labelFinally := "label" + strconv.Itoa(ce.labelCount)
	err = ce.vmwriter.WriteIf(labelFinally)
	if err != nil {
		return err
	}
	ce.labelCount++

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

	// goto labelWhile
	err = ce.vmwriter.WriteGoto(labelWhile)
	if err != nil {
		return err
	}
	// labelFinally
	err = ce.vmwriter.WriteLabel(labelFinally)
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileDo() error { // process the do keyword
	err := ce.ProcessKeyWord(tk.DO)
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
	// remove the return value from the stack
	ce.vmwriter.WritePop(vw.TEMP, 0)
	return nil
}

func (ce *CompilationEngine) CompileReturn() error {
	// process the return keyword
	err := ce.ProcessKeyWord(tk.RETURN)
	if err != nil {
		return err
	}

	// process return; for void functions
	if ce.t.CurrentToken.Val == tk.SEMICOLON.Val {
		// TODO: activate this error check
		// if !ce.subroutineST.IsCurrentVoidFunc() {
		// 	return fmt.Errorf("bare return in a non-void function: %s", ce.subroutineST.CurrentScope.Name)
		// }

		// void function must push constant 0 to the stack
		err := ce.vmwriter.WritePush(vw.CONSTANT, 0)
		if err != nil {
			return err
		}
	} else if ce.t.CurrentToken.Val != tk.SEMICOLON.Val {
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
	// write the return command to the VM writer
	err = ce.vmwriter.WriteReturn()
	if err != nil {
		return err
	}
	return nil
}
