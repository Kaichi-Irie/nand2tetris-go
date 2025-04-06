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

func CreateCEwithFirstToken(xmlFile io.Writer, r io.Reader) *CompilationEngine {
	t, err := tokenizer.CreateTokenizerWithFirstToken(r)
	if err != nil {
		panic(err)
	}
	return &CompilationEngine{
		xmlFile: xmlFile,
		t:       t,
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

VARDEC:
	// class var dec
	for {
		token := ce.t.CurrentToken
		switch token {
		case tokenizer.KeywordsMap[tokenizer.KT_STATIC]:
			err = ce.CompileClassVarDec(tokenizer.KT_STATIC)
			if err != nil {
				return err
			}
		case tokenizer.KeywordsMap[tokenizer.KT_FIELD]:
			err = ce.CompileClassVarDec(tokenizer.KT_FIELD)
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
		switch token {
		case tokenizer.KeywordsMap[tokenizer.KT_CONSTRUCTOR], tokenizer.KeywordsMap[tokenizer.KT_FUNCTION], tokenizer.KeywordsMap[tokenizer.KT_METHOD]:
			err = ce.CompileSubroutine()
			if err != nil {
				return err
			}
		default:
			break SUBROUTINEDEC
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
	case token == tokenizer.KeywordsMap[tokenizer.KT_CONSTRUCTOR]:
		err = ce.t.ProcessKeyWord(tokenizer.KT_CONSTRUCTOR, ce.xmlFile)
		if err != nil {
			return err
		}
	case token == tokenizer.KeywordsMap[tokenizer.KT_FUNCTION]:
		err = ce.t.ProcessKeyWord(tokenizer.KT_FUNCTION, ce.xmlFile)
		if err != nil {
			return err
		}
	case token == tokenizer.KeywordsMap[tokenizer.KT_METHOD]:
		err = ce.t.ProcessKeyWord(tokenizer.KT_METHOD, ce.xmlFile)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected token %s", token)
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
	err = ce.CompileStatements()
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

func (ce *CompilationEngine) CompileStatements() error {
	_, err := io.WriteString(ce.xmlFile, "<statements>\n")
	if err != nil {
		return err
	}

	// process the statements
LOOP:
	for {
		switch token := ce.t.CurrentToken; {
		case token == tokenizer.KeywordsMap[tokenizer.KT_LET]:
			err = ce.CompileLet()
			if err != nil {
				return err
			}
		case token == tokenizer.KeywordsMap[tokenizer.KT_IF]:
			err = ce.CompileIf()
			if err != nil {
				return err
			}
		case token == tokenizer.KeywordsMap[tokenizer.KT_WHILE]:
			err = ce.CompileWhile()
			if err != nil {
				return err
			}
		case token == tokenizer.KeywordsMap[tokenizer.KT_DO]:
			err = ce.CompileDo()
			if err != nil {
				return err
			}
		case token == tokenizer.KeywordsMap[tokenizer.KT_RETURN]:
			err = ce.CompileReturn()
			if err != nil {
				return err
			}
		default:
			break LOOP
		}
	}

	_, err = io.WriteString(ce.xmlFile, "</statements>\n")
	if err != nil {
		return err
	}
	return nil
}

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
	err = ce.CompileExpression(false)
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

func (ce *CompilationEngine) CompileIf() error {
	_, err := io.WriteString(ce.xmlFile, "<ifStatement>\n")
	if err != nil {
		return err
	}
	// process the if keyword
	err = ce.t.ProcessKeyWord(tokenizer.KT_IF, ce.xmlFile)
	if err != nil {
		return err
	}
	// process the (
	err = ce.t.ProcessSymbol("(", ce.xmlFile)
	if err != nil {
		return err
	}
	// process the expression
	err = ce.CompileExpression(false)
	if err != nil {
		return err
	}
	// process the )
	err = ce.t.ProcessSymbol(")", ce.xmlFile)
	if err != nil {
		return err
	}
	// process the {
	err = ce.t.ProcessSymbol("{", ce.xmlFile)
	if err != nil {
		return err
	}
	// process the statements
	err = ce.CompileStatements()
	if err != nil {
		return err
	}
	// process the }
	err = ce.t.ProcessSymbol("}", ce.xmlFile)
	if err != nil {
		return err
	}
	// process the else keyword
	if ce.t.CurrentToken == tokenizer.KeywordsMap[tokenizer.KT_ELSE] {
		err = ce.t.ProcessKeyWord(tokenizer.KT_ELSE, ce.xmlFile)
		if err != nil {
			return err
		}
		// process the {
		err = ce.t.ProcessSymbol("{", ce.xmlFile)
		if err != nil {
			return err
		}
		// process the statements
		err = ce.CompileStatements()
		if err != nil {
			return err
		}
		// process the }
		err = ce.t.ProcessSymbol("}", ce.xmlFile)
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(ce.xmlFile, "</ifStatement>\n")
	if err != nil {
		return err
	}
	return nil

}

func (ce *CompilationEngine) CompileWhile() error {
	_, err := io.WriteString(ce.xmlFile, "<whileStatement>\n")
	if err != nil {
		return err
	}

	// process the while keyword
	err = ce.t.ProcessKeyWord(tokenizer.KT_WHILE, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the (
	err = ce.t.ProcessSymbol("(", ce.xmlFile)
	if err != nil {
		return err
	}

	// process the expression
	err = ce.CompileExpression(false)
	if err != nil {
		return err
	}

	// process the )
	err = ce.t.ProcessSymbol(")", ce.xmlFile)
	if err != nil {
		return err
	}

	// process the {
	err = ce.t.ProcessSymbol("{", ce.xmlFile)
	if err != nil {
		return err
	}

	// process the statements
	err = ce.CompileStatements()
	if err != nil {
		return err
	}

	// process the }
	err = ce.t.ProcessSymbol("}", ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</whileStatement>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileDo() error {
	_, err := io.WriteString(ce.xmlFile, "<doStatement>\n")
	if err != nil {
		return err
	}

	// process the do keyword
	err = ce.t.ProcessKeyWord(tokenizer.KT_DO, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the subroutine call. Skip the <term> and <expression> tags
	err = ce.CompileExpression(true)
	if err != nil {
		return err
	}

	// process the ;
	err = ce.t.ProcessSymbol(";", ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</doStatement>\n")
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) CompileReturn() error {
	_, err := io.WriteString(ce.xmlFile, "<returnStatement>\n")
	if err != nil {
		return err
	}

	// process the return keyword
	err = ce.t.ProcessKeyWord(tokenizer.KT_RETURN, ce.xmlFile)
	if err != nil {
		return err
	}

	// process the expression. Skip if the case is return;
	if ce.t.CurrentToken != ";" {
		// process the expression
		err = ce.CompileExpression(false)
		if err != nil {
			return err
		}
	}
	// process the ;
	err = ce.t.ProcessSymbol(";", ce.xmlFile)
	if err != nil {
		return err
	}

	_, err = io.WriteString(ce.xmlFile, "</returnStatement>\n")
	if err != nil {
		return err
	}
	return nil
}

/*
CompileTerm compiles a term and writes it to the XML file.
Term: integerConstant | stringConstant | keywordConstant | varName | varName '[' expression ']' | subroutineCall | '(' expression ')' | unaryOp term
skipTags: if true, do not write the <term> and </term> tags. This is used for subroutine calls.
*/
func (ce *CompilationEngine) CompileTerm(isSubroutineCall bool) error {
	var err error

	/* process the subroutine call:
	subroutineCall: (className | varName) '.' subroutineName '(' expressionList ') or subroutineName '(' expressionList ')'
	*/
	if isSubroutineCall {
		// process the subroutine name or class name or var name
		err = ce.t.ProcessIdentifier(ce.xmlFile)
		if err != nil {
			return err
		}
		// process the . or (
		if ce.t.CurrentToken == "." {
			err = ce.t.ProcessSymbol(".", ce.xmlFile)
			if err != nil {
				return err
			}
			// process the subroutine name
			err = ce.t.ProcessIdentifier(ce.xmlFile)
			if err != nil {
				return err
			}
		}

		// process the (
		err = ce.t.ProcessSymbol("(", ce.xmlFile)
		if err != nil {
			return err
		}

		// process the expression list
		err = ce.CompileExpressionList()
		if err != nil {
			return err
		}
		// process the )
		err = ce.t.ProcessSymbol(")", ce.xmlFile)
		if err != nil {
			return err
		}
		return nil
	}

	// process not a subroutine call
	_, err = io.WriteString(ce.xmlFile, "<term>\n")
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

/*
CompileExpression compiles an expression and writes it to the XML file.
Expression: term (op term)*
op: + - * / & | < > =
skipTags: if true, do not write the <expression> and </expression> tags. This is used for subroutine calls.
*/
func (ce *CompilationEngine) CompileExpression(isSubroutineCall bool) error {
	var err error
	// skip the <expression> and <term> tags if this is a subroutine call
	if isSubroutineCall {
		return ce.CompileTerm(isSubroutineCall)
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
	// TODO: implement the operator

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
	if ce.t.CurrentToken == ")" {
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
	for ce.t.CurrentToken == "," {
		err = ce.t.ProcessSymbol(",", ce.xmlFile)
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
