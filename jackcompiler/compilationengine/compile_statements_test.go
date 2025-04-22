package compilationengine

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompileLet(t *testing.T) {

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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
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
	}
}

func TestCompileIf(t *testing.T) {
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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
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
	}
}

func TestCompileWhile(t *testing.T) {
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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
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
	}
}

func TestCompileDo(t *testing.T) {
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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
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
	}
}

func TestCompileReturn(t *testing.T) {
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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode))
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
	}
}
