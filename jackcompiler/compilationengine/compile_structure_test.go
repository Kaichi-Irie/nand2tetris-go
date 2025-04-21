package compilationengine

import (
	"bytes"
	tk "nand2tetris-go/jackcompiler/tokenizer"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompileClass(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := CompilationEngine{
		writer: xmlFile,
		t:      tk.New(strings.NewReader("")),
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
`}, {
			jackCode: `// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/10/ExpressionLessSquare/Main.jack


class Main {
    static boolean test;    // Added for testing -- there is no static keyword
                            // in the Square files.

    function void main() {
        var SquareGame game;
        let game = game;
        do game.run();
        do game.dispose();
        return;
    }

    function void more() {  // Added to test Jack syntax that is not used in
        var boolean b;      // the Square files.
        if (b) {
        }
        else {              // There is no else keyword in the Square files.
        }
        return;
    }
}
`,
			expectedXML: `<class>
<keyword> class </keyword>
<identifier> Main </identifier>
<symbol> { </symbol>
<classVarDec>
<keyword> static </keyword>
<keyword> boolean </keyword>
<identifier> test </identifier>
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
</subroutineDec>
<subroutineDec>
<keyword> function </keyword>
<keyword> void </keyword>
<identifier> more </identifier>
<symbol> ( </symbol>
<parameterList>
</parameterList>
<symbol> ) </symbol>
<subroutineBody>
<symbol> { </symbol>
<varDec>
<keyword> var </keyword>
<keyword> boolean </keyword>
<identifier> b </identifier>
<symbol> ; </symbol>
</varDec>
<statements>
<ifStatement>
<keyword> if </keyword>
<symbol> ( </symbol>
<expression>
<term>
<identifier> b </identifier>
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
		ce.t, err = tk.NewWithFirstToken(strings.NewReader(test.jackCode))
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
		fieldOrStatic tk.Token
		expectedXML   string
	}{
		{
			jackCode:      `static int i;`,
			fieldOrStatic: tk.STATIC,
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
			fieldOrStatic: tk.FIELD,
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
			fieldOrStatic: tk.STATIC,
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
			fieldOrStatic: tk.FIELD,
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
		ce.t, err = tk.NewWithFirstToken(strings.NewReader(test.jackCode))
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
		ce.t, _ = tk.NewWithFirstToken(strings.NewReader(test.jackCode))
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
		ce.t, _ = tk.NewWithFirstToken(strings.NewReader(test.jackCode))
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
		ce.t, _ = tk.NewWithFirstToken(strings.NewReader(test.jackCode))
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
		ce.t, _ = tk.NewWithFirstToken(strings.NewReader(test.jackCode))
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
