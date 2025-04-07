package compilationengine

import (
	"fmt"
	"io"
	"nand2tetris-go/jackcompiler/tokenizer"
)

/*
CompileClass compiles a class and writes it to the XML file.
Class: 'class' className '{' classVarDec* subroutineDec* '}'
*/
func (ce *CompilationEngine) CompileClass() error {
	io.WriteString(ce.xmlFile, "<class>\n")

	// class keyword
	err := ce.t.ProcessKeyWord(tokenizer.CLASS, ce.xmlFile)

	if err != nil {
		return err
	}

	// class name
	err = ce.t.ProcessIdentifier(ce.xmlFile)
	if err != nil {
		return err
	}

	// {
	err = ce.t.ProcessSymbol(tokenizer.LBRACE, ce.xmlFile)
	if err != nil {
		return err
	}

VARDEC:
	// class var dec
	for {
		token := ce.t.CurrentToken
		switch token.Val() {
		case tokenizer.STATIC.Val():
			err = ce.CompileClassVarDec(tokenizer.STATIC)
			if err != nil {
				return err
			}
		case tokenizer.FIELD.Val():
			err = ce.CompileClassVarDec(tokenizer.FIELD)
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
		switch token.Val() {
		case tokenizer.CONSTRUCTOR.Val(), tokenizer.FUNCTION.Val(), tokenizer.METHOD.Val():
			err = ce.CompileSubroutine()
			if err != nil {
				return err
			}
		default:
			break SUBROUTINEDEC
		}
	}
	// }
	err = ce.t.ProcessSymbol(tokenizer.RBRACE, ce.xmlFile)
	if err != nil {
		return err
	}

	io.WriteString(ce.xmlFile, "</class>\n")
	return nil
}

/*
CompileClassVarDec compiles a class variable declaration and writes it to the XML file.
ClassVarDec: (static | field) type varName (',' varName)* ';'
*/
func (ce *CompilationEngine) CompileClassVarDec(staticOrField tokenizer.Token) error {
	_, err := io.WriteString(ce.xmlFile, "<classVarDec>\n")
	if err != nil {
		return err
	}

	// static or field keyword
	err = ce.t.ProcessKeyWord(staticOrField, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the type: int, char, boolean, className
	err = ce.t.ProcessType(ce.xmlFile)
	if err != nil {
		return err
	}

	// process the var name
	err = ce.t.ProcessIdentifier(ce.xmlFile)
	if err != nil {
		return err
	}

	// process the comma or semicolon
	for ce.t.CurrentToken.Val() == tokenizer.COMMA.Val() {
		// process the comma
		err = ce.t.ProcessSymbol(tokenizer.COMMA, ce.xmlFile)
		if err != nil {
			return err
		}
		// process the var name
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}
	}
	// process the semicolon
	err = ce.t.ProcessSymbol(tokenizer.SEMICOLON, ce.xmlFile)
	if err != nil {
		return err
	}
	_, err = io.WriteString(ce.xmlFile, "</classVarDec>\n")
	if err != nil {
		return err
	}
	return nil

}

func (ce *CompilationEngine) CompileSubroutine() error {
	_, err := io.WriteString(ce.xmlFile, "<subroutineDec>\n")
	if err != nil {
		return err
	}

	// subroutine keyword: constructor, function, method
	switch token := ce.t.CurrentToken; {
	case token.Val() == tokenizer.CONSTRUCTOR.Val():
		err = ce.t.ProcessKeyWord(tokenizer.CONSTRUCTOR, ce.xmlFile)
		if err != nil {
			return err
		}
	case token.Val() == tokenizer.FUNCTION.Val():
		err = ce.t.ProcessKeyWord(tokenizer.FUNCTION, ce.xmlFile)
		if err != nil {
			return err
		}
	case token.Val() == tokenizer.METHOD.Val():
		err = ce.t.ProcessKeyWord(tokenizer.METHOD, ce.xmlFile)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected token %s", token.Val())
	}

	// process the void or type: int, char, boolean, className
	if token := ce.t.CurrentToken; token.Val() == "void" {
		err = ce.t.ProcessKeyWord(tokenizer.VOID, ce.xmlFile)
		if err != nil {
			return err
		}
	} else {
		// int, char, boolean, className
		err = ce.t.ProcessType(ce.xmlFile)
		if err != nil {
			return err
		}
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

	// process the parameter list
	err = ce.CompileParameterList()
	if err != nil {
		return err
	}

	// process the )
	err = ce.t.ProcessSymbol(tokenizer.RPAREN, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the subroutine body
	err = ce.CompileSubroutineBody()
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</subroutineDec>\n")
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
	_, err := io.WriteString(ce.xmlFile, "<varDec>\n")
	if err != nil {
		return err
	}

	// var keyword
	err = ce.t.ProcessKeyWord(tokenizer.VAR, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the type: int, char, boolean, className
	err = ce.t.ProcessType(ce.xmlFile)
	if err != nil {
		return err
	}

	// process the var name
	err = ce.t.ProcessIdentifier(ce.xmlFile)
	if err != nil {
		return err
	}
	// process the comma or semicolon
	for ce.t.CurrentToken.Val() == tokenizer.COMMA.Val() {
		// process the comma
		err = ce.t.ProcessSymbol(tokenizer.COMMA, ce.xmlFile)
		if err != nil {
			return err
		}
		// process the var name
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}
	}
	// process the semicolon
	err = ce.t.ProcessSymbol(tokenizer.SEMICOLON, ce.xmlFile)
	if err != nil {
		return err
	}
	_, err = io.WriteString(ce.xmlFile, "</varDec>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileParameterList() error {

	// no parameters
	if ce.t.CurrentToken.Val() == tokenizer.RPAREN.Val() {
		_, err := io.WriteString(ce.xmlFile, "<parameterList>\n</parameterList>\n")
		if err != nil {
			return err
		}
		return nil
	}

	// parameter list
	_, err := io.WriteString(ce.xmlFile, "<parameterList>\n")
	if err != nil {
		return err
	}

	// process the type: int, char, boolean, className
	err = ce.t.ProcessType(ce.xmlFile)
	if err != nil {
		return err
	}

	// process the var name
	err = ce.t.ProcessIdentifier(ce.xmlFile)
	if err != nil {
		return err
	}

	// process the comma or semicolon
	// process the comma

	for ce.t.CurrentToken.Val() == tokenizer.COMMA.Val() {
		err = ce.t.ProcessSymbol(tokenizer.COMMA, ce.xmlFile)
		if err != nil {
			return err
		}
		// process the type: int, char, boolean, className
		err = ce.t.ProcessType(ce.xmlFile)
		if err != nil {
			return err
		}
		// process the var name
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}

	}

	_, err = io.WriteString(ce.xmlFile, "</parameterList>\n")
	if err != nil {
		return err
	}

	return nil
}

func (ce *CompilationEngine) CompileSubroutineBody() error {
	_, err := io.WriteString(ce.xmlFile, "<subroutineBody>\n")
	if err != nil {
		return err
	}

	// process the {
	err = ce.t.ProcessSymbol(tokenizer.LBRACE, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the var dec
	for ce.t.CurrentToken.Val() == tokenizer.VAR.Val() {
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
	err = ce.t.ProcessSymbol(tokenizer.RBRACE, ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</subroutineBody>\n")
	if err != nil {
		return err
	}
	return nil
}
