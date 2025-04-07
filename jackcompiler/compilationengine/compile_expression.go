package compilationengine

import (
	"fmt"
	"io"
	"nand2tetris-go/jackcompiler/tokenizer"
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
		_, err = io.WriteString(ce.xmlFile, "<term>\n")
		if err != nil {
			return err
		}
	}

	// process the term
	switch token := ce.t.CurrentToken; {

	// process the '(' expression ')'
	case token == tokenizer.LPAREN:
		err = ce.t.ProcessSymbol(tokenizer.LPAREN, ce.xmlFile)
		if err != nil {
			return err
		}
		err = ce.CompileExpression(false)
		if err != nil {
			return err
		}
		err = ce.t.ProcessSymbol(tokenizer.RPAREN, ce.xmlFile)
		if err != nil {
			return err
		}

	/*
		process the unary operator, unaryOp: - ~
	*/
	case token.IsUnaryOp():
		err = ce.t.ProcessSymbol(token, ce.xmlFile)
		if err != nil {
			return err
		}
		err = ce.CompileTerm(false)
		if err != nil {
			return err
		}
	// process the string constant
	case token.Is(tokenizer.TT_STRING_CONST):
		err = ce.t.ProcessStringConst(ce.xmlFile)
		if err != nil {
			return err
		}
	// process the integer constant
	case token.Is(tokenizer.TT_INT_CONST):
		err = ce.t.ProcessIntConst(ce.xmlFile)
		if err != nil {
			return err
		}
	// process the keyword constant: true, false, null, this
	case token.IsKeywordConst():
		err = ce.t.ProcessKeyWord(token, ce.xmlFile)
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
	case token.Is(tokenizer.TT_IDENTIFIER):
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}
		// process the . or ( or [
		switch ce.t.CurrentToken.Val() {
		case tokenizer.DOT.Val():
			err = ce.t.ProcessSymbol(tokenizer.DOT, ce.xmlFile)
			if err != nil {
				return err
			}
			// process the subroutine name
			err = ce.t.ProcessIdentifier(ce.xmlFile)
			if err != nil {
				return err
			}
			// process the (
			err = ce.t.ProcessSymbol(tokenizer.LPAREN, ce.xmlFile)
			if err != nil {
				return err
			}
			// process the expression list
			err = ce.CompileExpressionList()
			if err != nil {
				return err
			}
			// process the )
			err = ce.t.ProcessSymbol(tokenizer.RPAREN, ce.xmlFile)
			if err != nil {
				return err
			}
		case tokenizer.LPAREN.Val():
			err = ce.t.ProcessSymbol(tokenizer.LPAREN, ce.xmlFile)
			if err != nil {
				return err
			}
			// process the expression list
			err = ce.CompileExpressionList()
			if err != nil {
				return err
			}
			// process the )
			err = ce.t.ProcessSymbol(tokenizer.RPAREN, ce.xmlFile)
			if err != nil {
				return err
			}
		case tokenizer.LSQUARE.Val():
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
		// write the identifier to the xml file

	default:
		return fmt.Errorf("unexpected token %s", token.Val())
	}
	// Do Statement: skip the </term> tag
	if isDoStatement {
		return nil
	}
	_, err = io.WriteString(ce.xmlFile, "</term>\n")
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

	_, err = io.WriteString(ce.xmlFile, "<expression>\n")
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
		err = ce.t.ProcessSymbol(token, ce.xmlFile)
		if err != nil {
			return err
		}
		// process the term
		err = ce.CompileTerm(false)
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(ce.xmlFile, "</expression>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileExpressionList() error {
	_, err := io.WriteString(ce.xmlFile, "<expressionList>\n")
	if err != nil {
		return err
	}
	// no expressions
	if ce.t.CurrentToken.Val() == tokenizer.RPAREN.Val() {
		_, err := io.WriteString(ce.xmlFile, "</expressionList>\n")
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
	for ce.t.CurrentToken.Val() == tokenizer.COMMA.Val() {
		err = ce.t.ProcessSymbol(tokenizer.COMMA, ce.xmlFile)
		if err != nil {
			return err
		}
		err = ce.CompileExpression(false)
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(ce.xmlFile, "</expressionList>\n")
	if err != nil {
		return err
	}
	return nil
}
