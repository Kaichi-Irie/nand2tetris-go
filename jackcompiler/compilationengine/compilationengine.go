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

// CompileClass compiles a class and writes it to the XML file
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

// CompileClassVarDec compiles a class variable declaration and writes it to the XML file
func (ce *CompilationEngine) CompileClassVarDec(staticOrField tokenizer.KeyWordType) error {
	// static or field keyword
	err := ce.t.ProcessKeyWord(staticOrField, ce.xmlFile)
	if err != nil {
		return err
	}

	token := ce.t.CurrentToken
	// process the type
	// case 1: className
	if tokenizer.IsIdentifier(token) {
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}
		return nil
	}

	// case 2: keyword: int, char, boolean
	kwt, err := tokenizer.GetKeyWordType(token)
	if err != nil {
		return err
	} else if kwt != tokenizer.KT_INT && kwt != tokenizer.KT_CHAR && kwt != tokenizer.KT_BOOLEAN {
		return fmt.Errorf("expected int, char or boolean, got %d", kwt)
	}
	err = ce.t.ProcessKeyWord(kwt, ce.xmlFile)
	if err != nil {
		return err
	}

	// TODO: process the multiple var names
	// process the var name
	err = ce.t.ProcessIdentifier(ce.xmlFile)
	if err != nil {
		return err
	}

	// process the comma or semicolon
	for ; token == ","; token = ce.t.CurrentToken {
		// process the comma
		err = ce.t.ProcessSymbol(",", ce.xmlFile)
		if err != nil {
			return err
		}

		ce.t.ProcessIdentifier(ce.xmlFile)
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
