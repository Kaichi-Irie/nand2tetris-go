package compilationengine

import (
	"fmt"
	"strconv"

	st "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/symboltable"
	tk "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/tokenizer"
	vw "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/vmwriter"
)

/*
CompileTerm compiles a term and writes it to the XML file.
Term: integerConstant | stringConstant | keywordConstant | varName | varName '[' expression ']' | subroutineCall | '(' expression ')' | unaryOp term
skipTags: if true, do not write the <term> and </term> tags. This is used for subroutine calls.
*/
func (ce *CompilationEngine) CompileTerm(isDoStatement bool) error {
	var err error
	// process the term
	switch token := ce.t.CurrentToken; {

	// process the '(' expression ')'
	case token == tk.LPAREN:
		err = ce.ProcessSymbol(tk.LPAREN)
		if err != nil {
			return err
		}
		err = ce.CompileExpression(false)
		if err != nil {
			return err
		}
		err = ce.ProcessSymbol(tk.RPAREN)
		if err != nil {
			return err
		}

	/*
		process the unary operator, unaryOp: - ~
	*/
	case token.IsUnaryOp():
		err = ce.ProcessSymbol(token)
		if err != nil {
			return err
		}
		err = ce.CompileTerm(false)
		if err != nil {
			return err
		}
		vmCommand := map[string]string{tk.MINUS.Val: vw.NEG, tk.NOT.Val: vw.NOT}[token.Val]
		err = ce.vmwriter.WriteArithmetic(vmCommand)
		if err != nil {
			return err
		}

	// process the string constant
	case token.Is(tk.TT_STRING_CONST):
		trimmedStrConst := token.Val[1 : len(token.Val)-1]
		length := len(trimmedStrConst)
		err = ce.vmwriter.WritePush(vw.CONSTANT, length)
		if err != nil {
			return err
		}
		err = ce.vmwriter.WriteCall("String.new", 1)
		if err != nil {
			return err
		}
		for _, c := range trimmedStrConst {
			err = ce.vmwriter.WritePush(vw.CONSTANT, int(c))
			if err != nil {
				return err
			}

			err = ce.vmwriter.WriteCall("String.appendChar", 2)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}

		err = ce.ProcessStringConst()
		if err != nil {
			return err
		}
	// process the integer constant
	case token.Is(tk.TT_INT_CONST):
		n, err := strconv.Atoi(token.Val)
		if err != nil {
			return err
		}
		err = ce.vmwriter.WritePush(vw.CONSTANT, n)
		if err != nil {
			return err
		}
		err = ce.ProcessIntConst()
		if err != nil {
			return err
		}
	// process the keyword constant: true, false, null, this
	case token.Val == tk.NULL.Val || token.Val == tk.FALSE.Val:
		err = ce.vmwriter.WritePush(vw.CONSTANT, 0)
		if err != nil {
			return err
		}
		err = ce.ProcessKeyWord(token)
		if err != nil {
			return err
		}
	case token.Val == tk.TRUE.Val:
		err = ce.vmwriter.WritePush(vw.CONSTANT, 1)
		if err != nil {
			return err
		}
		err = ce.vmwriter.WriteArithmetic(vw.NEG)
		if err != nil {
			return err
		}
		err = ce.ProcessKeyWord(token)
		if err != nil {
			return err
		}
	case token.Val == tk.THIS.Val:
		err = ce.vmwriter.WritePush(vw.POINTER, 0)
		if err != nil {
			return err
		}
		err = ce.ProcessKeyWord(token)
		if err != nil {
			return err
		}
	/*
		process varName | varName '[' expression ']' | subroutineCall
		varName: identifier
		subroutineName: identifier
		className: identifier
		subroutineCall: (className | varName) '.' subroutineName '(' expressionList ') or subroutineName '(' expressionList ')'
		As for subroutineCall, we have to remove XML tags, <term> ,</term> , <expression> , and </expression>.
		We have to process the identifier first, then check if it is a subroutine call or a varName. To do this, we use the strings.Builder as a temporary buffer.
	*/
	case token.Is(tk.TT_IDENTIFIER):
		name1 := token.Val // used only for the subroutine call
		err = ce.ProcessIdentifier()
		if err != nil {
			return err
		}

		// process the . or ( or [
		switch token := ce.t.CurrentToken; token.Val {
		// process (className | varName) '.' subroutineName '(' expressionList ')
		// name1 '.' name2 '(' expressionList ')
		case tk.DOT.Val:
			err = ce.ProcessSymbol(tk.DOT)
			if err != nil {
				return err
			}

			// process the subroutine name
			name2 := ce.t.CurrentToken.Val
			err = ce.ProcessIdentifier()
			if err != nil {
				return err
			}
			// process the (
			err = ce.ProcessSymbol(tk.LPAREN)
			if err != nil {
				return err
			}

			var subroutineName string
			var nArgs int
			/*
				process 2 cases;
				case1. varName.methodName(...)
				case2. className.functionName(...)
			*/
			if id, ok := ce.Lookup(name1); ok && (id.Kind == st.FIELD || id.Kind == st.STATIC || id.Kind == st.ARG || id.Kind == st.VAR) {
				// case 1.
				className := id.T
				methodName := name2
				subroutineName = className + "." + methodName

				seg := vw.SegmentOfKind[id.Kind]
				index := id.Index
				err = ce.vmwriter.WritePush(seg, index)
				if err != nil {
					return err
				}
				nArgs, err = ce.CompileExpressionList()
				if err != nil {
					return err
				}
				nArgs++ // for variable pointer
			} else {
				// case 2.
				className := name1
				methodName := name2
				subroutineName = className + "." + methodName
				nArgs, err = ce.CompileExpressionList()
				if err != nil {
					return err
				}
			}

			// process the )
			err = ce.ProcessSymbol(tk.RPAREN)
			if err != nil {
				return err
			}
			// write vmcommand

			err = ce.vmwriter.WriteCall(subroutineName, nArgs)
			if err != nil {
				return err
			}

		// process methodName '(' expressionList ')'
		case tk.LPAREN.Val:
			err = ce.ProcessSymbol(tk.LPAREN)
			if err != nil {
				return err
			}

			// push the object
			err = ce.vmwriter.WritePush(vw.POINTER, 0)
			if err != nil {
				return err
			}

			// process the expression list
			nArgs, err := ce.CompileExpressionList()
			if err != nil {
				return err
			}
			nArgs++
			// process the )
			err = ce.ProcessSymbol(tk.RPAREN)
			if err != nil {
				return err
			}
			// write vmcommand
			className := ce.classST.CurrentScope.Name
			methodName := name1
			subroutineName := className + "." + methodName
			err = ce.vmwriter.WriteCall(subroutineName, nArgs)
			if err != nil {
				return err
			}

		// process the varName '[' expression ']' | varName
		case tk.LSQUARE.Val:
			arrayName := name1

			// TODO: make this process as a function; WritePush(varName)
			id, _ := ce.Lookup(arrayName)
			// TODO: implement this error
			// if !ok {
			// 	return fmt.Errorf("array %s is not defined. Term cannot be used", arrayName)
			// }
			seg := vw.SegmentOfKind[id.Kind]
			index := id.Index
			ce.vmwriter.WritePush(seg, index)
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

			err = ce.vmwriter.WriteArithmetic(vw.ADD)
			if err != nil {
				return err
			}
			err = ce.vmwriter.WritePop(vw.POINTER, 1)
			if err != nil {
				return err
			}
			err = ce.vmwriter.WritePush(vw.THAT, 0)
			if err != nil {
				return err
			}

		default:
			// process the variable
			if id, ok := ce.Lookup(name1); ok {
				seg := vw.SegmentOfKind[id.Kind]
				index := id.Index
				err = ce.vmwriter.WritePush(seg, index)
				if err != nil {
					return err
				}
			}
			// TODO: activate this error
			//  else {
			// 	return fmt.Errorf("variable %s is not defined. Term cannot be used", name1)
			// }
		}

	default:
		return fmt.Errorf("unexpected token %s", token.Val)
	}
	return nil
}

// /*
// CompileExpression compiles an expression and writes it to the XML file.
// Expression: term (op term)*
// op: + - * / & | < > =
// skipTags: if true, do not write the <expression> and </expression> tags. This is used for subroutine calls.
// */
func (ce *CompilationEngine) CompileExpression(isDoStatement bool) error {
	var err error
	// skip the <expression> and <term> tags if this is a subroutine call
	if isDoStatement {
		return ce.CompileTerm(isDoStatement)
	}
	// process the term
	err = ce.CompileTerm(false)
	if err != nil {
		return err
	}

	// process the operator
	for token := ce.t.CurrentToken; token.IsOp(); token = ce.t.CurrentToken {
		vmCommand := map[string]string{
			tk.PLUS.Val:     vw.ADD,
			tk.MINUS.Val:    vw.SUB,
			tk.ASTERISK.Val: vw.MUL,
			tk.SLASH.Val:    vw.DIV,
			tk.AND.Val:      vw.AND,
			tk.OR.Val:       vw.OR,
			tk.LESS.Val:     vw.LT,
			tk.GREATER.Val:  vw.GT,
			tk.EQUAL.Val:    vw.EQ,
		}[token.Val]
		err = ce.ProcessSymbol(token)
		if err != nil {
			return err
		}
		// process the term
		err = ce.CompileTerm(false)
		if err != nil {
			return err
		}
		err = ce.vmwriter.WriteArithmetic(vmCommand)
		if err != nil {
			return err
		}

	}
	return nil
}

func (ce *CompilationEngine) CompileExpressionList() (int, error) {
	nArgs := 0
	// no expressions
	if ce.t.CurrentToken.Val == tk.RPAREN.Val {
		return nArgs, nil
	}

	// process the expression
	err := ce.CompileExpression(false)
	if err != nil {
		return 0, err
	}
	nArgs++

	// process the comma
	for ce.t.CurrentToken.Val == tk.COMMA.Val {
		err = ce.ProcessSymbol(tk.COMMA)
		if err != nil {
			return 0, err
		}
		err = ce.CompileExpression(false)
		if err != nil {
			return 0, err
		}

		nArgs++
	}
	return nArgs, nil
}
