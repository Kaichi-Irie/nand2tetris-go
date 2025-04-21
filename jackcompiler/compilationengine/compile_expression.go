package compilationengine

import (
	"fmt"
	"io"
	tk "nand2tetris-go/jackcompiler/tokenizer"
)

/*
CompileTerm compiles a term and writes it to the XML file.
Term: integerConstant | stringConstant | keywordConstant | varName | varName '[' expression ']' | subroutineCall | '(' expression ')' | unaryOp term
skipTags: if true, do not write the <term> and </term> tags. This is used for subroutine calls.
*/
func (ce *CompilationEngine) CompileTerm(isDoStatement bool) error {
	var err error
	// process not a subroutine call
	if !isDoStatement {
		_, err = io.WriteString(ce.writer, "<term>\n")
		if err != nil {
			return err
		}
	}

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
	// process the string constant
	case token.Is(tk.TT_STRING_CONST):
		err = ce.ProcessStringConst()
		if err != nil {
			return err
		}
	// process the integer constant
	case token.Is(tk.TT_INT_CONST):
		err = ce.ProcessIntConst()
		if err != nil {
			return err
		}
	// process the keyword constant: true, false, null, this
	case token.IsKeywordConst():
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
		err = ce.ProcessIdentifier()
		if err != nil {
			return err
		}
		// process the . or ( or [
		switch ce.t.CurrentToken.Val {
		case tk.DOT.Val:
			err = ce.ProcessSymbol(tk.DOT)
			if err != nil {
				return err
			}
			// process the subroutine name
			err = ce.ProcessIdentifier()
			if err != nil {
				return err
			}
			// process the (
			err = ce.ProcessSymbol(tk.LPAREN)
			if err != nil {
				return err
			}
			// process the expression list
			err = ce.CompileExpressionList()
			if err != nil {
				return err
			}
			// process the )
			err = ce.ProcessSymbol(tk.RPAREN)
			if err != nil {
				return err
			}
		case tk.LPAREN.Val:
			err = ce.ProcessSymbol(tk.LPAREN)
			if err != nil {
				return err
			}
			// process the expression list
			err = ce.CompileExpressionList()
			if err != nil {
				return err
			}
			// process the )
			err = ce.ProcessSymbol(tk.RPAREN)
			if err != nil {
				return err
			}
		case tk.LSQUARE.Val:
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
		// write the identifier to the xml file

	default:
		return fmt.Errorf("unexpected token %s", token.Val)
	}
	// Do Statement: skip the </term> tag
	if isDoStatement {
		return nil
	}
	_, err = io.WriteString(ce.writer, "</term>\n")
	if err != nil {
		return err
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

	_, err = io.WriteString(ce.writer, "<expression>\n")
	if err != nil {
		return err
	}

	// process the term
	err = ce.CompileTerm(false)
	if err != nil {
		return err
	}

	// process the operator
	for token := ce.t.CurrentToken; token.IsOp(); token = ce.t.CurrentToken {
		err = ce.ProcessSymbol(token)
		if err != nil {
			return err
		}
		// process the term
		err = ce.CompileTerm(false)
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(ce.writer, "</expression>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileExpressionList() error {
	_, err := io.WriteString(ce.writer, "<expressionList>\n")
	if err != nil {
		return err
	}
	// no expressions
	if ce.t.CurrentToken.Val == tk.RPAREN.Val {
		_, err := io.WriteString(ce.writer, "</expressionList>\n")
		if err != nil {
			return err
		}
		return nil
	}

	// process the expression
	err = ce.CompileExpression(false)
	if err != nil {
		return err
	}

	// process the comma
	for ce.t.CurrentToken.Val == tk.COMMA.Val {
		err = ce.ProcessSymbol(tk.COMMA)
		if err != nil {
			return err
		}
		err = ce.CompileExpression(false)
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(ce.writer, "</expressionList>\n")
	if err != nil {
		return err
	}
	return nil
}
