package compilationengine

import (
	"bytes"
	"nand2tetris-go/jackcompiler/symboltable"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompileTerm(t *testing.T) {

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
		}, {
			jackCode: `arr[1]`,
			expectedXML: `<term>
<identifier> arr </identifier>
<symbol> [ </symbol>
<expression>
<term>
<integerConstant> 1 </integerConstant>
</term>
</expression>
<symbol> ] </symbol>
</term>
`},
		{
			jackCode: `arr[i]`,
			expectedXML: `<term>
<identifier> arr </identifier>
<symbol> [ </symbol>
<expression>
<term>
<identifier> i </identifier>
</term>
</expression>
<symbol> ] </symbol>
</term>
`,
		}, {
			jackCode: `-i`,
			expectedXML: `<term>
<symbol> - </symbol>
<term>
<identifier> i </identifier>
</term>
</term>
`,
		}, {
			jackCode: `~b`,
			expectedXML: `<term>
<symbol> ~ </symbol>
<term>
<identifier> b </identifier>
</term>
</term>
`,
		}, {
			jackCode: `(1)`,
			expectedXML: `<term>
<symbol> ( </symbol>
<expression>
<term>
<integerConstant> 1 </integerConstant>
</term>
</expression>
<symbol> ) </symbol>
</term>
`,
		}, {
			jackCode: `sub(1)`,
			expectedXML: `<term>
<identifier> sub </identifier>
<symbol> ( </symbol>
<expressionList>
<expression>
<term>
<integerConstant> 1 </integerConstant>
</term>
</expression>
</expressionList>
<symbol> ) </symbol>
</term>
`,
		}, {
			jackCode: `classVar.sub(1)`,
			expectedXML: `<term>
<identifier> classVar </identifier>
<symbol> . </symbol>
<identifier> sub </identifier>
<symbol> ( </symbol>
<expressionList>
<expression>
<term>
<integerConstant> 1 </integerConstant>
</term>
</expression>
</expressionList>
<symbol> ) </symbol>
</term>
`,
		},
	}
	for _, test := range tests {
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode), "")
		ce.classST = &symboltable.SymbolTable{}
		ce.subroutineST = &symboltable.SymbolTable{}
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
	}
}

func TestCompileExpression(t *testing.T) {

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
		xmlFile := &bytes.Buffer{}
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode), "")
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
	}
}

func TestCompileIntConst(t *testing.T) {

	tests := []struct {
		jackCode string
		vmCode   string
	}{
		{
			jackCode: `1+2`,
			vmCode: `push constant 1
push constant 2
add
`},
		{
			jackCode: `1-2`,
			vmCode: `push constant 1
push constant 2
sub
`},
		{
			jackCode: `1*2`,
			vmCode: `push constant 1
push constant 2
call Math.multiply 2
`},
		{jackCode: `1/2`,
			vmCode: `push constant 1
push constant 2
call Math.divide 2
`},
		{
			jackCode: `1&2`,
			vmCode: `push constant 1
push constant 2
and
`},
		{jackCode: `-1`,
			vmCode: `push constant 1
neg
`},
		{
			jackCode: `~1`,
			vmCode: `push constant 1
not
`},
	}
	for _, test := range tests {
		vmFile := &bytes.Buffer{}
		ce := NewWithVMWriter(vmFile, &bytes.Buffer{}, strings.NewReader(test.jackCode), "")
		err := ce.CompileExpression(false)
		if err != nil {
			t.Errorf("CompileExpression() error: %v", err)
		}
		if vmFile.String() != test.vmCode {
			t.Errorf("CompileExpression() = %v, want %v", vmFile.String(), test.vmCode)
			diff := cmp.Diff(vmFile.String(), test.vmCode)
			t.Errorf("Diff: %s", diff)
		}
	}
}
