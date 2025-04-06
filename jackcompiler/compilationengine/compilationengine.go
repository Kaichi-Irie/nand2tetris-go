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

	var typeToken string
	// process the type
	switch {
	case ce.t.CurrentToken == "int" || ce.t.CurrentToken == "char" || ce.t.CurrentToken == "boolean":
		typeToken = ce.t.CurrentToken // int, char, boolean
		err = ce.t.ProcessKeyWord(tokenizer.KT_INT, ce.xmlFile)
		if err != nil {
			return err
		}

	case tokenizer.IsIdentifier(ce.t.CurrentToken):
		typeToken = ce.t.CurrentToken // className
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected int, char, boolean or className, got %s", ce.t.CurrentToken)
	}

	// TODO: process the multiple var names
	// process the var name
	err = ce.t.ProcessIdentifier(ce.xmlFile)
	if err != nil {
		return err
	}

	// process the comma or semicolon
	for ce.t.CurrentToken == "," {
		// process the comma
		ce.t.CurrentToken = ";"
		err = ce.t.ProcessSymbol(";", ce.xmlFile)
		if err != nil {
			return err
		}

		_, err = io.WriteString(ce.xmlFile, "<keyword> "+tokenizer.KeywordsMap[staticOrField]+" </keyword>\n")
		if err != nil {
			return err
		}
		if typeToken == "int" || typeToken == "char" || typeToken == "boolean" {
			_, err = io.WriteString(ce.xmlFile, "<keyword> "+typeToken+" </keyword>\n")
			if err != nil {
				return err
			}

		} else {
			_, err = io.WriteString(ce.xmlFile, "<identifier> "+typeToken+" </identifier>\n")
			if err != nil {
				return err
			}

		}

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
	io.WriteString(ce.xmlFile, "<subroutineDec>\n")

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
	return nil
}

// func (ce *CompilationEngine) CompileParameterList() error
// func (ce *CompilationEngine) CompileSubroutineBody() error
// func (ce *CompilationEngine) CompileVarDec() error
// func (ce *CompilationEngine) CompileStatements() error
// func (ce *CompilationEngine) CompileLet() error
// func (ce *CompilationEngine) CompileIf() error
// func (ce *CompilationEngine) CompileWhile() error
// func (ce *CompilationEngine) CompileDo() error
// func (ce *CompilationEngine) CompileReturn() error
// func (ce *CompilationEngine) CompileExpression() error
// func (ce *CompilationEngine) CompileTerm() error
// func (ce *CompilationEngine) CompileExpressionList() error
