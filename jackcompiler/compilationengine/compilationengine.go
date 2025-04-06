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

	// process the type: int, char, boolean or className
	switch {
	case ce.t.CurrentToken == "int" || ce.t.CurrentToken == "char" || ce.t.CurrentToken == "boolean":
		err = ce.t.ProcessKeyWord(tokenizer.KT_INT, ce.xmlFile)
		if err != nil {
			return err
		}
	case tokenizer.IsIdentifier(ce.t.CurrentToken):
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected int, char, boolean or className, got %s", ce.t.CurrentToken)
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

	// subroutine keyword
	kwt, err := tokenizer.GetKeyWordType(ce.t.CurrentToken)
	if err != nil {
		return err
	} else if kwt != tokenizer.KT_CONSTRUCTOR && kwt != tokenizer.KT_FUNCTION && kwt != tokenizer.KT_METHOD {
		return fmt.Errorf("expected constructor, function or method, got %d", kwt)
	}

	err = ce.t.ProcessKeyWord(kwt, ce.xmlFile)
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

	// process the type: int, char, boolean or className
	switch {
	case ce.t.CurrentToken == "int" || ce.t.CurrentToken == "char" || ce.t.CurrentToken == "boolean":
		err = ce.t.ProcessKeyWord(tokenizer.KT_INT, ce.xmlFile)
		if err != nil {
			return err
		}
	case tokenizer.IsIdentifier(ce.t.CurrentToken):
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("expected int, char, boolean or className, got %s", ce.t.CurrentToken)
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

	// process the type: int, char, boolean or className
	switch token := ce.t.CurrentToken; {
	case token == "int" || token == "char" || token == "boolean":
		kwt, err := tokenizer.GetKeyWordType(token)
		if err != nil {
			return err
		}

		err = ce.t.ProcessKeyWord(kwt, ce.xmlFile)
		if err != nil {
			return err
		}
	case tokenizer.IsIdentifier(ce.t.CurrentToken):
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected int, char, boolean or className, got %s", ce.t.CurrentToken)
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

		// process the type: int, char, boolean or className
		switch token := ce.t.CurrentToken; {
		case token == "int" || token == "char" || token == "boolean":
			kwt, err := tokenizer.GetKeyWordType(token)
			if err != nil {
				return err
			}
			err = ce.t.ProcessKeyWord(kwt, ce.xmlFile)
			if err != nil {
				return err
			}

		case tokenizer.IsIdentifier(ce.t.CurrentToken):
			err = ce.t.ProcessIdentifier(ce.xmlFile)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("expected int, char, boolean or className, got %s", ce.t.CurrentToken)
		}

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

// func (ce *CompilationEngine) CompileSubroutineBody() error
// func (ce *CompilationEngine) CompileStatements() error
// func (ce *CompilationEngine) CompileLet() error
// func (ce *CompilationEngine) CompileIf() error
// func (ce *CompilationEngine) CompileWhile() error
// func (ce *CompilationEngine) CompileDo() error
// func (ce *CompilationEngine) CompileReturn() error
// func (ce *CompilationEngine) CompileExpression() error
// func (ce *CompilationEngine) CompileTerm() error
// func (ce *CompilationEngine) CompileExpressionList() error
