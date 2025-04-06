package compilationengine

import (
	"bytes"
	"nand2tetris-go/jackcompiler/tokenizer"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompileClass(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := CompilationEngine{
		xmlFile: xmlFile,
		t:       tokenizer.New(strings.NewReader("")),
	}

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `class Main {static int i;static int j;}`,
			expectedXML: `<class>
<keyword> class </keyword>
<identifier> Main </identifier>
<symbol> { </symbol>
<classVarDec>
<keyword> static </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> ; </symbol>
</classVarDec>
<classVarDec>
<keyword> static </keyword>
<keyword> int </keyword>
<identifier> j </identifier>
<symbol> ; </symbol>
</classVarDec>
<symbol> } </symbol>
</class>
`,
		}, {
			jackCode: `class Main {
			field MyClass i,j;
			}`,
			expectedXML: `<class>
<keyword> class </keyword>
<identifier> Main </identifier>
<symbol> { </symbol>
<classVarDec>
<keyword> field </keyword>
<identifier> MyClass </identifier>
<identifier> i </identifier>
<symbol> , </symbol>
<identifier> j </identifier>
<symbol> ; </symbol>
</classVarDec>
<symbol> } </symbol>
</class>
`,
		}, {
			jackCode: `class Main {
			field MyClass i,j;
			function void main() {return;}
			}`,
			expectedXML: `<class>
<keyword> class </keyword>
<identifier> Main </identifier>
<symbol> { </symbol>
<classVarDec>
<keyword> field </keyword>
<identifier> MyClass </identifier>
<identifier> i </identifier>
<symbol> , </symbol>
<identifier> j </identifier>
<symbol> ; </symbol>
</classVarDec>
<subroutineDec>
<keyword> function </keyword>
<keyword> void </keyword>
<identifier> main </identifier>
<symbol> ( </symbol>
<parameterList>
</parameterList>
<symbol> ) </symbol>
<subroutineBody>
<symbol> { </symbol>
<statements>
<returnStatement>
<keyword> return </keyword>
<symbol> ; </symbol>
</returnStatement>
</statements>
<symbol> } </symbol>
</subroutineBody>
</subroutineDec>
<symbol> } </symbol>
</class>
`},
	}
	var err error
	for _, test := range tests {
		ce.t, err = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		if err != nil {
			t.Errorf("CreateTokenizerWithFirstToken() error: %v", err)
		}
		err = ce.CompileClass()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}

}

func TestCompileClassVarDec(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode      string
		fieldOrStatic tokenizer.KeyWordType
		expectedXML   string
	}{
		{
			jackCode:      `static int i;`,
			fieldOrStatic: tokenizer.KT_STATIC,
			expectedXML: `<classVarDec>
<keyword> static </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> ; </symbol>
</classVarDec>
`,
		},
		{
			jackCode:      `field int i,j;`,
			fieldOrStatic: tokenizer.KT_FIELD,
			expectedXML: `<classVarDec>
<keyword> field </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> , </symbol>
<identifier> j </identifier>
<symbol> ; </symbol>
</classVarDec>
`,
		},
		{
			jackCode:      `static int i,j;`,
			fieldOrStatic: tokenizer.KT_STATIC,
			expectedXML: `<classVarDec>
<keyword> static </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> , </symbol>
<identifier> j </identifier>
<symbol> ; </symbol>
</classVarDec>
`,
		},
		{
			jackCode:      `field MyClass A,B,C;`,
			fieldOrStatic: tokenizer.KT_FIELD,
			expectedXML: `<classVarDec>
<keyword> field </keyword>
<identifier> MyClass </identifier>
<identifier> A </identifier>
<symbol> , </symbol>
<identifier> B </identifier>
<symbol> , </symbol>
<identifier> C </identifier>
<symbol> ; </symbol>
</classVarDec>
`,
		},
	}
	var err error
	// show diff using cmp package
	for _, test := range tests {
		ce.t, err = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		if err != nil {
			t.Errorf("CreateTokenizerWithFirstToken() error: %v", err)
		}
		err = ce.CompileClassVarDec(test.fieldOrStatic)
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}

		xmlFile.Reset()
	}

}

func TestCompileVarDec(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `var int i,j;`,
			expectedXML: `<varDec>
<keyword> var </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> , </symbol>
<identifier> j </identifier>
<symbol> ; </symbol>
</varDec>
`,
		},
		{
			jackCode: `var MyClass i;`,
			expectedXML: `<varDec>
<keyword> var </keyword>
<identifier> MyClass </identifier>
<identifier> i </identifier>
<symbol> ; </symbol>
</varDec>
`,
		}}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileVarDec()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestCompileParameterList(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `int i, boolean b, MyClass c`,
			expectedXML: `<parameterList>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> , </symbol>
<keyword> boolean </keyword>
<identifier> b </identifier>
<symbol> , </symbol>
<identifier> MyClass </identifier>
<identifier> c </identifier>
</parameterList>
`,
		},
		{
			jackCode: `int i`,
			expectedXML: `<parameterList>
<keyword> int </keyword>
<identifier> i </identifier>
</parameterList>
`,
		},
	}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileParameterList()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestCompileTerm(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `i`,
			expectedXML: `<term>
<identifier> i </identifier>
</term>
`,
		}, {
			jackCode: `123`,
			expectedXML: `<term>
<integerConstant> 123 </integerConstant>
</term>
`,
		}, {
			jackCode: `true`,
			expectedXML: `<term>
<keyword> true </keyword>
</term>
`,
		}, {
			jackCode: "\"hello world\"",
			expectedXML: `<term>
<stringConstant> hello world </stringConstant>
</term>
`,
		}, {
			jackCode: `this`,
			expectedXML: `<term>
<keyword> this </keyword>
</term>
`,
		}, {
			jackCode: `false`,
			expectedXML: `<term>
<keyword> false </keyword>
</term>
`,
		}, {
			jackCode: `null`,
			expectedXML: `<term>
<keyword> null </keyword>
</term>
`,
		}}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileTerm(false)
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestCompileExpression(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `i`,
			expectedXML: `<expression>
<term>
<identifier> i </identifier>
</term>
</expression>
`,
		}, {
			jackCode: "\"hello world\"",
			expectedXML: `<expression>
<term>
<stringConstant> hello world </stringConstant>
</term>
</expression>
`,
		}}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileExpression(false)
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestCompileLet(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `let i = j;`,
			expectedXML: `<letStatement>
<keyword> let </keyword>
<identifier> i </identifier>
<symbol> = </symbol>
<expression>
<term>
<identifier> j </identifier>
</term>
</expression>
<symbol> ; </symbol>
</letStatement>
`,
		},
	}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileLet()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestCompileReturn(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `return i;`,
			expectedXML: `<returnStatement>
<keyword> return </keyword>
<expression>
<term>
<identifier> i </identifier>
</term>
</expression>
<symbol> ; </symbol>
</returnStatement>
`,
		},
		{
			jackCode: `return;`,
			expectedXML: `<returnStatement>
<keyword> return </keyword>
<symbol> ; </symbol>
</returnStatement>
`},
	}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileReturn()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestCompileWhile(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `while (i) { }`,
			expectedXML: `<whileStatement>
<keyword> while </keyword>
<symbol> ( </symbol>
<expression>
<term>
<identifier> i </identifier>
</term>
</expression>
<symbol> ) </symbol>
<symbol> { </symbol>
<statements>
</statements>
<symbol> } </symbol>
</whileStatement>
`}, {
			jackCode: `while (i) { let i = i; }`,
			expectedXML: `<whileStatement>
<keyword> while </keyword>
<symbol> ( </symbol>
<expression>
<term>
<identifier> i </identifier>
</term>
</expression>
<symbol> ) </symbol>
<symbol> { </symbol>
<statements>
<letStatement>
<keyword> let </keyword>
<identifier> i </identifier>
<symbol> = </symbol>
<expression>
<term>
<identifier> i </identifier>
</term>
</expression>
<symbol> ; </symbol>
</letStatement>
</statements>
<symbol> } </symbol>
</whileStatement>
`},
	}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileWhile()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestCompileDo(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `do myfunc();`,
			expectedXML: `<doStatement>
<keyword> do </keyword>
<identifier> myfunc </identifier>
<symbol> ( </symbol>
<expressionList>
</expressionList>
<symbol> ) </symbol>
<symbol> ; </symbol>
</doStatement>
`}, {
			jackCode: `do game.run();`,
			expectedXML: `<doStatement>
<keyword> do </keyword>
<identifier> game </identifier>
<symbol> . </symbol>
<identifier> run </identifier>
<symbol> ( </symbol>
<expressionList>
</expressionList>
<symbol> ) </symbol>
<symbol> ; </symbol>
</doStatement>
`},
	}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileDo()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestCompileExpressionList(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `i, j, k`,
			expectedXML: `<expressionList>
<expression>
<term>
<identifier> i </identifier>
</term>
</expression>
<symbol> , </symbol>
<expression>
<term>
<identifier> j </identifier>
</term>
</expression>
<symbol> , </symbol>
<expression>
<term>
<identifier> k </identifier>
</term>
</expression>
</expressionList>
`}}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileExpressionList()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestCompileIf(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `if (x) { }`,
			expectedXML: `<ifStatement>
<keyword> if </keyword>
<symbol> ( </symbol>
<expression>
<term>
<identifier> x </identifier>
</term>
</expression>
<symbol> ) </symbol>
<symbol> { </symbol>
<statements>
</statements>
<symbol> } </symbol>
</ifStatement>
`},
		{
			jackCode: `if (x) {  } else {  }`,
			expectedXML: `<ifStatement>
<keyword> if </keyword>
<symbol> ( </symbol>
<expression>
<term>
<identifier> x </identifier>
</term>
</expression>
<symbol> ) </symbol>
<symbol> { </symbol>
<statements>
</statements>
<symbol> } </symbol>
<keyword> else </keyword>
<symbol> { </symbol>
<statements>
</statements>
<symbol> } </symbol>
</ifStatement>
`}}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileIf()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}
func TestCompileSubroutineBody(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `{ }`,
			expectedXML: `<subroutineBody>
<symbol> { </symbol>
<statements>
</statements>
<symbol> } </symbol>
</subroutineBody>
`,
		}, {
			jackCode: `{ var int i; }`,
			expectedXML: `<subroutineBody>
<symbol> { </symbol>
<varDec>
<keyword> var </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> ; </symbol>
</varDec>
<statements>
</statements>
<symbol> } </symbol>
</subroutineBody>
`,
		}, {
			jackCode: `{ var SquareGame game;
let game = game;
do game.run();
do game.dispose();
return; }`,
			expectedXML: `<subroutineBody>
<symbol> { </symbol>
<varDec>
<keyword> var </keyword>
<identifier> SquareGame </identifier>
<identifier> game </identifier>
<symbol> ; </symbol>
</varDec>
<statements>
<letStatement>
<keyword> let </keyword>
<identifier> game </identifier>
<symbol> = </symbol>
<expression>
<term>
<identifier> game </identifier>
</term>
</expression>
<symbol> ; </symbol>
</letStatement>
<doStatement>
<keyword> do </keyword>
<identifier> game </identifier>
<symbol> . </symbol>
<identifier> run </identifier>
<symbol> ( </symbol>
<expressionList>
</expressionList>
<symbol> ) </symbol>
<symbol> ; </symbol>
</doStatement>
<doStatement>
<keyword> do </keyword>
<identifier> game </identifier>
<symbol> . </symbol>
<identifier> dispose </identifier>
<symbol> ( </symbol>
<expressionList>
</expressionList>
<symbol> ) </symbol>
<symbol> ; </symbol>
</doStatement>
<returnStatement>
<keyword> return </keyword>
<symbol> ; </symbol>
</returnStatement>
</statements>
<symbol> } </symbol>
</subroutineBody>
`,
		}}
	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileSubroutineBody()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}

func TestSubroutine(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `function void main() {
			var int i;
			return;
			}`,
			expectedXML: `<subroutineDec>
<keyword> function </keyword>
<keyword> void </keyword>
<identifier> main </identifier>
<symbol> ( </symbol>
<parameterList>
</parameterList>
<symbol> ) </symbol>
<subroutineBody>
<symbol> { </symbol>
<varDec>
<keyword> var </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> ; </symbol>
</varDec>
<statements>
<returnStatement>
<keyword> return </keyword>
<symbol> ; </symbol>
</returnStatement>
</statements>
<symbol> } </symbol>
</subroutineBody>
</subroutineDec>
`,
		}}

	for _, test := range tests {
		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		err := ce.CompileSubroutine()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
			t.Errorf("Diff: %s", diff)
		}
		xmlFile.Reset()
	}
}
