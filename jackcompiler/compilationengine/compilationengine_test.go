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
