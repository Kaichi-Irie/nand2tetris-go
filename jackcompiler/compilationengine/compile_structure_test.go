package compilationengine

import (
	"bytes"
	st "nand2tetris-go/jackcompiler/symboltable"
	tk "nand2tetris-go/jackcompiler/tokenizer"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompileClass(t *testing.T) {
	tests := []struct {
		jackCode             string
		expectedXML          string
		expectedClassST      st.SymbolTable
		expectedSubroutineST st.SymbolTable
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
			expectedClassST: st.SymbolTable{
				VarCnt:    0,
				FieldCnt:  0,
				StaticCnt: 2,
				ArgCnt:    0,
				VariableMap: map[string]st.Identifier{
					"Main": {
						Name:  "Main",
						Kind:  st.NONE,
						T:     "Main",
						Index: 0},
					"i": {
						Name:  "i",
						Kind:  tk.STATIC.Val,
						T:     tk.INT.Val,
						Index: 0},
					"j": {
						Name:  "j",
						Kind:  tk.STATIC.Val,
						T:     tk.INT.Val,
						Index: 1},
				},
			},
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
			expectedClassST: st.SymbolTable{
				VarCnt:    0,
				FieldCnt:  2,
				StaticCnt: 0,
				ArgCnt:    0,
				VariableMap: map[string]st.Identifier{
					"Main": {
						Name:  "Main",
						Kind:  st.NONE,
						T:     "Main",
						Index: 0},
					"i": {
						Name:  "i",
						Kind:  tk.FIELD.Val,
						T:     "MyClass",
						Index: 0},
					"j": {
						Name:  "j",
						Kind:  tk.FIELD.Val,
						T:     "MyClass",
						Index: 1},
				},
			},
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
`,
			expectedClassST: st.SymbolTable{
				VarCnt:    0,
				FieldCnt:  2,
				StaticCnt: 0,
				ArgCnt:    0,
				VariableMap: map[string]st.Identifier{
					"Main": {
						Name:  "Main",
						Kind:  st.NONE,
						T:     "Main",
						Index: 0},
					"i": {
						Name:  "i",
						Kind:  tk.FIELD.Val,
						T:     "MyClass",
						Index: 0},
					"j": {
						Name:  "j",
						Kind:  tk.FIELD.Val,
						T:     "MyClass",
						Index: 1},
					"main": {
						Name:  "main",
						Kind:  st.NONE,
						T:     "subroutine",
						Index: 0},
				},
			}},
		{
			jackCode: `class Main {
			field MyClass i,j;
			function void main(int a, boolean b) {
			var int c;
			var boolean d,e;
			return;}`,
			expectedClassST: st.SymbolTable{
				VarCnt:    0,
				FieldCnt:  2,
				StaticCnt: 0,
				ArgCnt:    0,
				VariableMap: map[string]st.Identifier{
					"Main": {
						Name:  "Main",
						Kind:  st.NONE,
						T:     "Main",
						Index: 0},
					"i": {
						Name:  "i",
						Kind:  tk.FIELD.Val,
						T:     "MyClass",
						Index: 0},
					"j": {
						Name:  "j",
						Kind:  tk.FIELD.Val,
						T:     "MyClass",
						Index: 1},
					"main": {
						Name:  "main",
						Kind:  st.NONE,
						T:     "subroutine",
						Index: 0},
				}}, expectedSubroutineST: st.SymbolTable{
				VarCnt:    3,
				FieldCnt:  0,
				StaticCnt: 0,
				ArgCnt:    2,
				VariableMap: map[string]st.Identifier{
					"a": {
						Name:  "a",
						Kind:  st.ARG,
						T:     "int",
						Index: 0},
					"b": {
						Name:  "b",
						Kind:  st.ARG,
						T:     "boolean",
						Index: 1},
					"c": {
						Name:  "c",
						Kind:  st.VAR,
						T:     "int",
						Index: 0},
					"d": {
						Name:  "d",
						Kind:  st.VAR,
						T:     "boolean",
						Index: 1},
					"e": {
						Name:  "e",
						Kind:  st.VAR,
						T:     "boolean",
						Index: 2},
				},
			},
		},

		{
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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
		err = ce.CompileClass()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if want := test.expectedXML; want != "" && xmlFile.String() != want {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), want)
			diff := cmp.Diff(xmlFile.String(), want)
			t.Errorf("Diff: %s", diff)
		}

		// compare the symbol table with go-cmp
		if want := test.expectedClassST; !cmp.Equal(want, st.SymbolTable{}) && !cmp.Equal(*ce.classST, want) {
			t.Errorf("CompileClass() = %v, want %v", ce.classST, want)
			diff := cmp.Diff(ce.classST, want)
			t.Errorf("Diff: %s", diff)
		}
		if want := test.expectedSubroutineST; !cmp.Equal(want, st.SymbolTable{}) && !cmp.Equal(*ce.subroutineST, want) {
			t.Errorf("CompileClass() = %v, want %v", ce.subroutineST, want)
			diff := cmp.Diff(ce.subroutineST, want)
			t.Errorf("Diff: %s", diff)
		}
		// reset the xmlFile and symbol table for the next test
	}

}

func TestCompileClassVarDec(t *testing.T) {

	tests := []struct {
		jackCode            string
		fieldOrStatic       tk.Token
		expectedXML         string
		expectedSymbolTable st.SymbolTable
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
			expectedSymbolTable: st.SymbolTable{
				VarCnt:    0,
				FieldCnt:  0,
				StaticCnt: 1,
				ArgCnt:    0,
				VariableMap: map[string]st.Identifier{
					"i": {
						Name:  "i",
						Kind:  tk.STATIC.Val,
						T:     tk.INT.Val,
						Index: 0},
				},
			},
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
			expectedSymbolTable: st.SymbolTable{
				VarCnt:    0,
				FieldCnt:  2,
				StaticCnt: 0,
				ArgCnt:    0,
				VariableMap: map[string]st.Identifier{
					"i": {
						Name:  "i",
						Kind:  tk.FIELD.Val,
						T:     tk.INT.Val,
						Index: 0},
					"j": {
						Name:  "j",
						Kind:  tk.FIELD.Val,
						T:     tk.INT.Val,
						Index: 1},
				},
			},
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
			expectedSymbolTable: st.SymbolTable{
				VarCnt:    0,
				FieldCnt:  0,
				StaticCnt: 2,
				ArgCnt:    0,
				VariableMap: map[string]st.Identifier{
					"i": {
						Name:  "i",
						Kind:  tk.STATIC.Val,
						T:     tk.INT.Val,
						Index: 0},
					"j": {
						Name:  "j",
						Kind:  tk.STATIC.Val,
						T:     tk.INT.Val,
						Index: 1},
				}}},
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
			expectedSymbolTable: st.SymbolTable{
				VarCnt:    0,
				FieldCnt:  3,
				StaticCnt: 0,
				ArgCnt:    0,
				VariableMap: map[string]st.Identifier{
					"A": {
						Name:  "A",
						Kind:  tk.FIELD.Val,
						T:     "MyClass",
						Index: 0},
					"B": {
						Name:  "B",
						Kind:  tk.FIELD.Val,
						T:     "MyClass",
						Index: 1},
					"C": {
						Name:  "C",
						Kind:  tk.FIELD.Val,
						T:     "MyClass",
						Index: 2},
				},
			}},
	}
	var err error
	// show diff using cmp package
	for _, test := range tests {
		xmlFile := &bytes.Buffer{}
		ce := New(xmlFile, strings.NewReader(""))
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

		// compare the symbol table with go-cmp
		if !cmp.Equal(*ce.classST, test.expectedSymbolTable) {
			t.Errorf("CompileClass() = %v, want %v", ce.classST, test.expectedSymbolTable)
			diff := cmp.Diff(ce.classST, test.expectedSymbolTable)
			t.Errorf("Diff: %s", diff)
		}

		xmlFile.Reset()
	}

}

func TestSubroutine(t *testing.T) {

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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
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
	}
}

func TestCompileVarDec(t *testing.T) {
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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
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
	}
}

func TestCompileParameterList(t *testing.T) {

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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
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
	}
}

func TestCompileSubroutineBody(t *testing.T) {
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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
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
	}
}
