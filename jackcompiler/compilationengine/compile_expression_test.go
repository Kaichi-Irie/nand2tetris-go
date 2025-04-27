package compilationengine

import (
	"bytes"

	st "nand2tetris-go/jackcompiler/symboltable"
	tk "nand2tetris-go/jackcompiler/tokenizer"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompileStringConst(t *testing.T) {
	tests := []struct {
		jackCode          string
		expectedVMCommand string
	}{
		{
			jackCode: `"hello"`,
			expectedVMCommand: `push constant 5
call String.new 1
push constant 104
call String.appendChar 2
push constant 101
call String.appendChar 2
push constant 108
call String.appendChar 2
push constant 108
call String.appendChar 2
push constant 111
call String.appendChar 2
`,
		},
	}
	for _, test := range tests {
		vmFile := &bytes.Buffer{}
		ce := NewWithVMWriter(vmFile, strings.NewReader(test.jackCode), "")
		err := ce.CompileTerm(false)
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		if vmFile.String() != test.expectedVMCommand {
			t.Errorf("CompileClass() = %v, want %v", vmFile.String(), test.expectedVMCommand)
			diff := cmp.Diff(vmFile.String(), test.expectedVMCommand)
			t.Errorf("Diff: %s", diff)
		}
	}
}

func TestArray(t *testing.T) {
	variables := map[string]st.Identifier{
		"arr": {
			Name:  "arr",
			Kind:  st.VAR,
			T:     st.ARRAY,
			Index: 0,
		},
		"i": {
			Name:  "i",
			Kind:  st.VAR,
			T:     tk.INT.Val,
			Index: 1,
		}}
	tests := []struct {
		jackCode          string
		expectedVMCommand string
	}{
		{
			jackCode: `arr[i]`,
			expectedVMCommand: `push local 0
push local 1
add
pop pointer 1
push that 0
`,
		}}
	for _, test := range tests {
		vmFile := &bytes.Buffer{}
		ce := NewWithVMWriter(vmFile, strings.NewReader(test.jackCode), "")
		ce.classST = &st.SymbolTable{}
		ce.subroutineST = &st.SymbolTable{}
		ce.subroutineST.VariableMap = variables
		err := ce.CompileTerm(false)
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		if vmFile.String() != test.expectedVMCommand {
			t.Errorf("CompileClass() = %v, want %v", vmFile.String(), test.expectedVMCommand)
			diff := cmp.Diff(vmFile.String(), test.expectedVMCommand)
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
		ce := NewWithVMWriter(vmFile, strings.NewReader(test.jackCode), "")
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
