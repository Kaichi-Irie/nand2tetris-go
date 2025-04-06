package compilationengine

import (
	"fmt"
	"io"
	"nand2tetris-go/jackcompiler/tokenizer"
)

type CompilationEngine struct {
	xmlFile io.Writer
	t       *tokenizer.Tokenizer
}

func New(xmlFile io.Writer, r io.Reader) *CompilationEngine {
	return &CompilationEngine{
		xmlFile: xmlFile,
		t:       tokenizer.New(r),
	}
}

/*
CompileClass compiles a class and writes it to the XML file.
Class: 'class' className '{' classVarDec* subroutineDec* '}'
*/
func (ce *CompilationEngine) CompileClass() error {
	io.WriteString(ce.xmlFile, "<class>\n")

	// class keyword
	err := ce.t.ProcessKeyWord(tokenizer.KT_CLASS, ce.xmlFile)

	if err != nil {
		return err
	}

	// class name
	err = ce.t.ProcessIdentifier(ce.xmlFile)
	if err != nil {
		return err
	}

	// {
	err = ce.t.ProcessSymbol("{", ce.xmlFile)
	if err != nil {
		return err
	}

	// class var dec
	for {
		token := ce.t.CurrentToken

		if kwt, _ := tokenizer.GetKeyWordType(token); kwt == tokenizer.KT_STATIC || kwt == tokenizer.KT_FIELD {
			err = ce.CompileClassVarDec(kwt)
			if err != nil {
				return err
			}
			continue
		} else {
			break
		}

	}

	// }
	err = ce.t.ProcessSymbol("}", ce.xmlFile)
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
func (ce *CompilationEngine) CompileClassVarDec(staticOrField tokenizer.KeyWordType) error {
	// static or field keyword
	err := ce.t.ProcessKeyWord(staticOrField, ce.xmlFile)
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
	for ce.t.CurrentToken == "," {
		// process the comma
		err = ce.t.ProcessSymbol(",", ce.xmlFile)
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
	err = ce.t.ProcessSymbol(";", ce.xmlFile)
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
	kwt, err := tokenizer.GetKeyWordType(ce.t.CurrentToken)
	if err != nil {
		return err
	} else if kwt != tokenizer.KT_CONSTRUCTOR && kwt != tokenizer.KT_FUNCTION && kwt != tokenizer.KT_METHOD {
		return fmt.Errorf("expected constructor, function or method, got %d", kwt)
	}

	// process the void or type: int, char, boolean, className
	if token := ce.t.CurrentToken; token == "void" {
		err = ce.t.ProcessKeyWord(tokenizer.KT_VOID, ce.xmlFile)
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
	err = ce.t.ProcessSymbol("(", ce.xmlFile)
	if err != nil {
		return err
	}

	// process the parameter list
	err = ce.CompileParameterList()
	if err != nil {
		return err
	}

	// process the )
	err = ce.t.ProcessSymbol(")", ce.xmlFile)
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
	err = ce.t.ProcessKeyWord(tokenizer.KT_VAR, ce.xmlFile)
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
	for ce.t.CurrentToken == "," {
		// process the comma
		err = ce.t.ProcessSymbol(",", ce.xmlFile)
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
	err = ce.t.ProcessSymbol(";", ce.xmlFile)
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
	if ce.t.CurrentToken == ")" {
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
	for ce.t.CurrentToken == "," {
		err = ce.t.ProcessSymbol(",", ce.xmlFile)
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
	err = ce.t.ProcessSymbol("{", ce.xmlFile)
	if err != nil {
		return err
	}

	// process the var dec
	for ce.t.CurrentToken == "var" {
		err = ce.CompileVarDec()
		if err != nil {
			return err
		}
	}

	// process the statements
	// TODO: implement the statements
	// err = ce.CompileStatements()
	if err != nil {
		return err
	}

	// process the }
	err = ce.t.ProcessSymbol("}", ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</subroutineBody>\n")
	if err != nil {
		return err
	}
	return nil
}

// func (ce *CompilationEngine) CompileStatements() error

func (ce *CompilationEngine) CompileLet() error {
	_, err := io.WriteString(ce.xmlFile, "<letStatement>\n")
	if err != nil {
		return err
	}

	// process the let keyword
	err = ce.t.ProcessKeyWord(tokenizer.KT_LET, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the var name
	err = ce.t.ProcessIdentifier(ce.xmlFile)
	if err != nil {
		return err
	}

	// TODO: process the [ or =

	// process the =
	err = ce.t.ProcessSymbol("=", ce.xmlFile)
	if err != nil {
		return err
	}

	// process the expression
	err = ce.CompileExpression()
	if err != nil {
		return err
	}

	// process the ;
	err = ce.t.ProcessSymbol(";", ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</letStatement>\n")
	if err != nil {
		return err
	}
	return nil
}

// func (ce *CompilationEngine) CompileIf() error
// func (ce *CompilationEngine) CompileWhile() error
// func (ce *CompilationEngine) CompileDo() error
// func (ce *CompilationEngine) CompileReturn() error
func (ce *CompilationEngine) CompileTerm() error {
	_, err := io.WriteString(ce.xmlFile, "<term>\n")
	if err != nil {
		return err
	}

	// process the term
	// TODO: implement the term
	switch token := ce.t.CurrentToken; {
	case tokenizer.IsIdentifier(token):
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}
	case tokenizer.IsStringConst(token):
		err = ce.t.ProcessStringConst(ce.xmlFile)
		if err != nil {
			return err
		}

	case tokenizer.IsIntConst(token):
		err = ce.t.ProcessIntConst(ce.xmlFile)
		if err != nil {
			return err
		}
	case token == tokenizer.KeywordsMap[tokenizer.KT_TRUE]:
		err = ce.t.ProcessKeyWord(tokenizer.KT_TRUE, ce.xmlFile)
		if err != nil {
			return err
		}
	case token == tokenizer.KeywordsMap[tokenizer.KT_FALSE]:
		err = ce.t.ProcessKeyWord(tokenizer.KT_FALSE, ce.xmlFile)
		if err != nil {
			return err
		}
	case token == tokenizer.KeywordsMap[tokenizer.KT_NULL]:
		err = ce.t.ProcessKeyWord(tokenizer.KT_NULL, ce.xmlFile)
		if err != nil {
			return err
		}
	case token == tokenizer.KeywordsMap[tokenizer.KT_THIS]:
		err = ce.t.ProcessKeyWord(tokenizer.KT_THIS, ce.xmlFile)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected token %s", token)
	}

	_, err = io.WriteString(ce.xmlFile, "</term>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileExpression() error {
	_, err := io.WriteString(ce.xmlFile, "<expression>\n")
	if err != nil {
		return err
	}

	// process the term
	err = ce.CompileTerm()
	if err != nil {
		return err
	}

	// process the operator
	// TODO: implement the operator

	_, err = io.WriteString(ce.xmlFile, "</expression>\n")
	if err != nil {
		return err
	}
	return nil
}

// func (ce *CompilationEngine) CompileExpressionList() error
