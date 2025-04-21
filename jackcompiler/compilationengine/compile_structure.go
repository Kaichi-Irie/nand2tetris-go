package compilationengine

import (
	"fmt"
	"io"
	tk "nand2tetris-go/jackcompiler/tokenizer"
)

/*
CompileClass compiles a class and writes it to the XML file.
Class: 'class' className '{' classVarDec* subroutineDec* '}'
*/
func (ce *CompilationEngine) CompileClass() error {
	io.WriteString(ce.writer, "<class>\n")

	// class keyword
	err := ce.ProcessKeyWord(tk.CLASS)

	if err != nil {
		return err
	}

	// class name
	err = ce.ProcessIdentifier()
	if err != nil {
		return err
	}

	// {
	err = ce.ProcessSymbol(tk.LBRACE)
	if err != nil {
		return err
	}

VARDEC:
	// class var dec
	for {
		token := ce.t.CurrentToken
		switch token.Val {
		case tk.STATIC.Val:
			err = ce.CompileClassVarDec(tk.STATIC)
			if err != nil {
				return err
			}
		case tk.FIELD.Val:
			err = ce.CompileClassVarDec(tk.FIELD)
			if err != nil {
				return err
			}
		default:
			break VARDEC
		}
	}

SUBROUTINEDEC:
	// subroutine dec
	for {
		token := ce.t.CurrentToken
		switch token.Val {
		case tk.CONSTRUCTOR.Val, tk.FUNCTION.Val, tk.METHOD.Val:
			err = ce.CompileSubroutine()
			if err != nil {
				return err
			}
		default:
			break SUBROUTINEDEC
		}
	}
	// }
	err = ce.ProcessSymbol(tk.RBRACE)
	if err != nil {
		return err
	}

	io.WriteString(ce.writer, "</class>\n")
	return nil
}

/*
CompileClassVarDec compiles a class variable declaration and writes it to the XML file.
ClassVarDec: (static | field) type varName (',' varName)* ';'
*/
func (ce *CompilationEngine) CompileClassVarDec(staticOrField tk.Token) error {
	_, err := io.WriteString(ce.writer, "<classVarDec>\n")
	if err != nil {
		return err
	}

	// static or field keyword
	err = ce.ProcessKeyWord(staticOrField)
	if err != nil {
		return err
	}

	// process the type: int, char, boolean, className
	err = ce.ProcessType()
	if err != nil {
		return err
	}

	// process the var name
	err = ce.ProcessIdentifier()
	if err != nil {
		return err
	}

	// process the comma or semicolon
	for ce.t.CurrentToken.Val == tk.COMMA.Val {
		// process the comma
		err = ce.ProcessSymbol(tk.COMMA)
		if err != nil {
			return err
		}
		// process the var name
		err = ce.ProcessIdentifier()
		if err != nil {
			return err
		}
	}
	// process the semicolon
	err = ce.ProcessSymbol(tk.SEMICOLON)
	if err != nil {
		return err
	}
	_, err = io.WriteString(ce.writer, "</classVarDec>\n")
	if err != nil {
		return err
	}
	return nil

}

func (ce *CompilationEngine) CompileSubroutine() error {
	_, err := io.WriteString(ce.writer, "<subroutineDec>\n")
	if err != nil {
		return err
	}

	// subroutine keyword: constructor, function, method
	switch token := ce.t.CurrentToken; {
	case token.Val == tk.CONSTRUCTOR.Val:
		err = ce.ProcessKeyWord(tk.CONSTRUCTOR)
		if err != nil {
			return err
		}
	case token.Val == tk.FUNCTION.Val:
		err = ce.ProcessKeyWord(tk.FUNCTION)
		if err != nil {
			return err
		}
	case token.Val == tk.METHOD.Val:
		err = ce.ProcessKeyWord(tk.METHOD)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected token %s", token.Val)
	}

	// process the void or type: int, char, boolean, className
	if token := ce.t.CurrentToken; token.Val == "void" {
		err = ce.ProcessKeyWord(tk.VOID)
		if err != nil {
			return err
		}
	} else {
		// int, char, boolean, className
		err = ce.ProcessType()
		if err != nil {
			return err
		}
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

	// process the parameter list
	err = ce.CompileParameterList()
	if err != nil {
		return err
	}

	// process the )
	err = ce.ProcessSymbol(tk.RPAREN)
	if err != nil {
		return err
	}

	// process the subroutine body
	err = ce.CompileSubroutineBody()
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.writer, "</subroutineDec>\n")
	if err != nil {
		return err
	}
	return nil
}

/*
CompileVarDec compiles a variable declaration and writes it to the XML file.
VarDec: 'var' type varName (',' varName)* ';'
*/
func (ce *CompilationEngine) CompileVarDec() error {
	_, err := io.WriteString(ce.writer, "<varDec>\n")
	if err != nil {
		return err
	}

	// var keyword
	err = ce.ProcessKeyWord(tk.VAR)
	if err != nil {
		return err
	}

	// process the type: int, char, boolean, className
	err = ce.ProcessType()
	if err != nil {
		return err
	}

	// process the var name
	err = ce.ProcessIdentifier()
	if err != nil {
		return err
	}
	// process the comma or semicolon
	for ce.t.CurrentToken.Val == tk.COMMA.Val {
		// process the comma
		err = ce.ProcessSymbol(tk.COMMA)
		if err != nil {
			return err
		}
		// process the var name
		err = ce.ProcessIdentifier()
		if err != nil {
			return err
		}
	}
	// process the semicolon
	err = ce.ProcessSymbol(tk.SEMICOLON)
	if err != nil {
		return err
	}
	_, err = io.WriteString(ce.writer, "</varDec>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileParameterList() error {

	// no parameters
	if ce.t.CurrentToken.Val == tk.RPAREN.Val {
		_, err := io.WriteString(ce.writer, "<parameterList>\n</parameterList>\n")
		if err != nil {
			return err
		}
		return nil
	}

	// parameter list
	_, err := io.WriteString(ce.writer, "<parameterList>\n")
	if err != nil {
		return err
	}

	// process the type: int, char, boolean, className
	err = ce.ProcessType()
	if err != nil {
		return err
	}

	// process the var name
	err = ce.ProcessIdentifier()
	if err != nil {
		return err
	}

	// process the comma or semicolon
	// process the comma

	for ce.t.CurrentToken.Val == tk.COMMA.Val {
		err = ce.ProcessSymbol(tk.COMMA)
		if err != nil {
			return err
		}
		// process the type: int, char, boolean, className
		err = ce.ProcessType()
		if err != nil {
			return err
		}
		// process the var name
		err = ce.ProcessIdentifier()
		if err != nil {
			return err
		}

	}

	_, err = io.WriteString(ce.writer, "</parameterList>\n")
	if err != nil {
		return err
	}

	return nil
}

func (ce *CompilationEngine) CompileSubroutineBody() error {
	_, err := io.WriteString(ce.writer, "<subroutineBody>\n")
	if err != nil {
		return err
	}

	// process the {
	err = ce.ProcessSymbol(tk.LBRACE)
	if err != nil {
		return err
	}

	// process the var dec
	for ce.t.CurrentToken.Val == tk.VAR.Val {
		err = ce.CompileVarDec()
		if err != nil {
			return err
		}
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

	_, err = io.WriteString(ce.writer, "</subroutineBody>\n")
	if err != nil {
		return err
	}
	return nil
}
