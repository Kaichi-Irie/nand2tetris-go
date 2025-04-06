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
	ce := New(xmlFile, strings.NewReader(""))

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
<keyword> static </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> ; </symbol>
<keyword> static </keyword>
<keyword> int </keyword>
<identifier> j </identifier>
<symbol> ; </symbol>
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
<keyword> field </keyword>
<identifier> MyClass </identifier>
<identifier> i </identifier>
<symbol> , </symbol>
<identifier> j </identifier>
<symbol> ; </symbol>
<symbol> } </symbol>
</class>
`,
		},
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
			expectedXML: `<keyword> static </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> ; </symbol>
`,
		},
		{
			jackCode:      `field int i,j;`,
			fieldOrStatic: tokenizer.KT_FIELD,
			expectedXML: `<keyword> field </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> , </symbol>
<identifier> j </identifier>
<symbol> ; </symbol>
`,
		},
		{
			jackCode:      `static int i,j;`,
			fieldOrStatic: tokenizer.KT_STATIC,
			expectedXML: `<keyword> static </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> , </symbol>
<identifier> j </identifier>
<symbol> ; </symbol>
`,
		},
		{
			jackCode:      `field MyClass A,B,C;`,
			fieldOrStatic: tokenizer.KT_FIELD,
			expectedXML: `<keyword> field </keyword>
<identifier> MyClass </identifier>
<identifier> A </identifier>
<symbol> , </symbol>
<identifier> B </identifier>
<symbol> , </symbol>
<identifier> C </identifier>
<symbol> ; </symbol>
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

// func TestCompileLet(t *testing.T) {
// 	xmlFile := &bytes.Buffer{}
// 	ce := New(xmlFile, strings.NewReader(""))

// 	tests := []struct {
// 		jackCode    string
// 		expectedXML string
// 	}{
// 		{
// 			jackCode: `let i = j;`,
// 			expectedXML: `<letStatement>
// <keyword> let </keyword>
// <identifier> i </identifier>
// <symbol> = </symbol>
// <identifier> j </identifier>
// <symbol> ; </symbol>
// </letStatement>
// `,
// 		},
// 	}
// 	for _, test := range tests {
// 		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
// 		err := ce.CompileLet()
// 		if err != nil {
// 			t.Errorf("CompileClass() error: %v", err)
// 		}
// 		// remove leading and trailing whitespace from the actual XML
// 		if xmlFile.String() != test.expectedXML {
// 			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
// 			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
// 			t.Errorf("Diff: %s", diff)
// 		}
// 		xmlFile.Reset()
// 	}
// }
// func TestCompileIf(t *testing.T) {
// 	xmlFile := &bytes.Buffer{}
// 	ce := New(xmlFile, strings.NewReader(""))

// 	tests := []struct {
// 		jackCode    string
// 		expectedXML string
// 	}{
// 		{
// 			jackCode: `if (x) { }`,
// 			expectedXML: `<ifStatement>
// <keyword> if </keyword>
// <symbol> ( </symbol>
// <expression>
// <term>
// <identifier> x </identifier>
// </term>
// </expression>
// <symbol> ) </symbol>
// <symbol> { </symbol>
// <symbol> } </symbol>
// </ifStatement>
// `},
// 		{
// 			jackCode: `if (x) {  } else {  }`,
// 			expectedXML: `<ifStatement>
// <keyword> if </keyword>
// <symbol> ( </symbol>
// <expression>
// <term>
// <identifier> x </identifier>
// </term>
// </expression>
// <symbol> ) </symbol>
// <symbol> { </symbol>
// <statements>
// </statements>
// <symbol> } </symbol>
// <keyword> else </keyword>
// <symbol> { </symbol>
// <statements>
// </statements>
// <symbol> } </symbol>
// </ifStatement>
// `}}
// 	for _, test := range tests {
// 		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
// 		err := ce.CompileIf()
// 		if err != nil {
// 			t.Errorf("CompileClass() error: %v", err)
// 		}
// 		// remove leading and trailing whitespace from the actual XML
// 		if xmlFile.String() != test.expectedXML {
// 			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
// 			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
// 			t.Errorf("Diff: %s", diff)
// 		}
// 		xmlFile.Reset()
// 	}
// }

// func TestSubroutine(t *testing.T) {
// 	xmlFile := &bytes.Buffer{}
// 	ce := New(xmlFile, strings.NewReader(""))

// 	tests := []struct {
// 		jackCode    string
// 		expectedXML string
// 	}{
// 		{
// 			jackCode: `function void main() {
// 			var int i;
// 			return;
// 			}`,
// 			expectedXML: `<subroutineDec>
// <keyword> function </keyword>
// <keyword> void </keyword>
// <identifier> main </identifier>
// <symbol> ( </symbol>
// <parameterList>
// </parameterList>
// <symbol> ) </symbol>
// <subroutineBody>
// <symbol> { </symbol>
// <varDec>
// <keyword> var </keyword>
// <keyword> int </keyword>
// <identifier> i </identifier>
// <symbol> ; </symbol>
// </varDec>
// <returnStatement>
// <keyword> return </keyword>
// <symbol> ; </symbol>
// </returnStatement>
// <symbol> } </symbol>
// </subroutineBody>
// </subroutineDec>
// `,
// 		}}

// 	for _, test := range tests {
// 		ce.t, _ = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
// 		err := ce.CompileSubroutine()
// 		if err != nil {
// 			t.Errorf("CompileClass() error: %v", err)
// 		}
// 		// remove leading and trailing whitespace from the actual XML
// 		if xmlFile.String() != test.expectedXML {
// 			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
// 			diff := cmp.Diff(xmlFile.String(), test.expectedXML)
// 			t.Errorf("Diff: %s", diff)
// 		}
// 		xmlFile.Reset()
// 	}
// }
