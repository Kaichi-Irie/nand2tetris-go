package compilationengine

import (
	"bytes"
	"strings"
	"testing"

	st "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/symboltable"
	tk "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/tokenizer"

	"github.com/google/go-cmp/cmp"
)

var codeComparer = cmp.Comparer(func(x string, y string) bool {
	x = strings.TrimSpace(x)
	y = strings.TrimSpace(y)
	x = strings.Join(strings.Fields(x), " ")
	y = strings.Join(strings.Fields(y), " ")
	return x == y
})

func TestCompileLet(t *testing.T) {
	tests := []struct {
		jackCode           string
		Variables          []st.Identifier
		expectedVMCommands string
	}{
		{
			jackCode: `let i = j;`,
			Variables: []st.Identifier{
				{Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0},
				{Name: "j", Kind: st.VAR, T: tk.INT.Val, Index: 1}},
			expectedVMCommands: `push local 1
pop local 0`,
		},
		{
			jackCode: `let i = j + k;`,
			Variables: []st.Identifier{
				{Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0},
				{Name: "j", Kind: st.VAR, T: tk.INT.Val, Index: 1},
				{Name: "k", Kind: st.VAR, T: tk.INT.Val, Index: 2}},
			expectedVMCommands: `push local 1
push local 2
add
pop local 0`,
		},
		{
			jackCode: `let i = j*k;`,
			Variables: []st.Identifier{
				{Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0},
				{Name: "j", Kind: st.VAR, T: tk.INT.Val, Index: 1},
				{Name: "k", Kind: st.VAR, T: tk.INT.Val, Index: 2}},
			expectedVMCommands: `push local 1
push local 2
call Math.multiply 2
pop local 0`,
		},
		{
			jackCode: `let variable = i + 1;`,
			Variables: []st.Identifier{
				{Name: "variable", Kind: st.STATIC, T: tk.INT.Val, Index: 0},
				{Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0}},
			expectedVMCommands: `push local 0
push constant 1
add
pop static 0`},
		{
			jackCode: `let b = false;`,
			Variables: []st.Identifier{
				{Name: "b", Kind: st.VAR, T: tk.BOOLEAN.Val, Index: 0}},
			expectedVMCommands: `push constant 0
pop local 0`,
		},
		{
			jackCode: `let b = true&false;`,
			Variables: []st.Identifier{
				{Name: "b", Kind: st.VAR, T: tk.BOOLEAN.Val, Index: 0}},
			expectedVMCommands: `push constant 1
neg
push constant 0
and
pop local 0`,
		},
		{
			jackCode: `let a[i] = b[j];`,
			Variables: []st.Identifier{
				{Name: "a", Kind: st.VAR, T: st.ARRAY, Index: 0},
				{Name: "b", Kind: st.VAR, T: st.ARRAY, Index: 1},
				{Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 2},
				{Name: "j", Kind: st.VAR, T: tk.INT.Val, Index: 3}},
			expectedVMCommands: `push local 0
push local 2
add
push local 1
push local 3
add
pop pointer 1
push that 0
pop temp 0
pop pointer 1
push temp 0
pop that 0`,
		},
	}

	for _, test := range tests {
		vmFile := &bytes.Buffer{}
		ce := NewWithVMWriter(vmFile, strings.NewReader(test.jackCode), "")

		// Define the variables in the symbol table
		for _, id := range test.Variables {
			ce.classST.Define(id.Name, id.T, id.Kind)
		}
		err := ce.CompileLet()
		if err != nil {
			t.Errorf("CompileLet() error: %v", err)
		}

		// trim leading and trailing whitespace
		vmOutput := vmFile.String()
		want := test.expectedVMCommands
		if diff := cmp.Diff(vmOutput, want, codeComparer); diff != "" {
			t.Errorf("CompileLet() = %v, want %v", vmOutput, want)
		}
	}
}

func TestCompileIf(t *testing.T) {
	Variables := []st.Identifier{
		{Name: "x", Kind: st.VAR, T: tk.INT.Val, Index: 0}}
	tests := []struct {
		jackCode          string
		expectedVMCommand string
	}{
		{
			jackCode: `if (x) { }`,
			expectedVMCommand: `push local 0
not
if-goto label0
goto label1
label label0
label label1`},
		{
			jackCode: `if (x) {  } else {  }`,
			expectedVMCommand: `push local 0
not
if-goto label0
goto label1
label label0
label label1`},
		{
			jackCode: `if (x) {let x=false;} else {let x=false;}`,
			expectedVMCommand: `push local 0
not
if-goto label0
push constant 0
pop local 0
goto label1
label label0
push constant 0
pop local 0
label label1`,
		},
	}
	for _, test := range tests {
		vmFile := &bytes.Buffer{}
		ce := NewWithVMWriter(vmFile, strings.NewReader(test.jackCode), "")

		// Define the variables in the symbol table
		for _, id := range Variables {
			ce.classST.Define(id.Name, id.T, id.Kind)
		}
		err := ce.CompileIf()
		if err != nil {
			t.Errorf("CompileLet() error: %v", err)
		}

		// trim leading and trailing whitespace
		vmOutput := vmFile.String()
		want := test.expectedVMCommand
		if diff := cmp.Diff(vmOutput, want, codeComparer); diff != "" {
			t.Errorf("CompileLet() = %v, want %v", vmOutput, want)
		}
	}
}

func TestCompileWhile(t *testing.T) {
	variables := []st.Identifier{
		{Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0}}
	tests := []struct {
		name              string
		jackCode          string
		expectedVMCommand string
	}{
		{
			name:     "Simple while",
			jackCode: `while (i) { }`,
			expectedVMCommand: `label label0
            push local 0
            not
            if-goto label1
            goto label0
            label label1`},
		{
			name:     "While with statement",
			jackCode: `while (i) { let i = i; }`,
			expectedVMCommand: `label label0
            push local 0
            not
            if-goto label1
            push local 0
            pop local 0
            goto label0
            label label1`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithVMWriter(vmFile, strings.NewReader(tt.jackCode), "TestClass")

			// Define the variables in the symbol table
			ce.subroutineST.SetCurrentScope("TestClass.testFunc", st.KINDFUNCTION, st.VOIDFUNC)
			for _, id := range variables {
				ce.subroutineST.Define(id.Name, id.T, id.Kind)
			}

			err := ce.CompileWhile()
			if err != nil {
				t.Fatalf("CompileWhile() error: %v", err)
			}

			vmOutput := vmFile.String()
			want := tt.expectedVMCommand
			if diff := cmp.Diff(want, vmOutput, codeComparer); diff != "" {
				t.Errorf("CompileWhile() VM code mismatch (-want +got):\n%s", diff)
				t.Logf("Got VM Code:\n%s", vmOutput)
				t.Logf("Want VM Code:\n%s", want)
			}
		})
	}
}

func TestCompileDo(t *testing.T) {
	tests := []struct {
		name               string
		jackCode           string
		className          string
		variables          []st.Identifier // For method calls on objects
		expectedVMCommands string
	}{
		{
			name:               "Simple function call",
			jackCode:           `do Output.printInt(1);`,
			className:          "Main",
			expectedVMCommands: `push constant 1 call Output.printInt 1 pop temp 0`,
		},
		{
			name:               "Method call on current object",
			jackCode:           `do draw();`, // Assumes draw is a method of the current class
			className:          "Square",
			expectedVMCommands: `push pointer 0 call Square.draw 1 pop temp 0`, // Pushes 'this', calls method, pops void return
		},
		{
			name:      "Method call on another object",
			jackCode:  `do game.run();`,
			className: "Main",
			variables: []st.Identifier{
				{Name: "game", Kind: st.VAR, T: "SquareGame", Index: 0}, // Assume 'game' is a local var of type SquareGame
			},
			expectedVMCommands: `push local 0 call SquareGame.run 1 pop temp 0`, // Pushes 'game' object, calls method, pops void return
		},
		{
			name:      "Function call with multiple args",
			jackCode:  `do Math.multiply(x, y);`,
			className: "Main",
			variables: []st.Identifier{ // Define x and y
				{Name: "x", Kind: st.VAR, T: tk.INT.Val, Index: 0},
				{Name: "y", Kind: st.VAR, T: tk.INT.Val, Index: 1},
			},
			expectedVMCommands: `push local 0 push local 1 call Math.multiply 2 pop temp 0`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithVMWriter(vmFile, strings.NewReader(tt.jackCode), tt.className)

			// Set up context
			ce.classST.SetCurrentScope(tt.className, st.KINDCLASS, st.NOTVOIDFUNC)
			ce.subroutineST.SetCurrentScope(tt.className+".testFunc", st.KINDFUNCTION, st.VOIDFUNC) // Assume a function context
			if strings.Contains(tt.jackCode, "draw()") {                                            // If it's a method call test
				ce.subroutineST.SetCurrentScope(tt.className+".testMethod", st.KINDMETHOD, st.VOIDFUNC)
				ce.subroutineST.Define("this", tt.className, st.ARG) // Add 'this' for method context
			}

			// Define variables if provided
			for _, id := range tt.variables {
				if id.Kind == st.STATIC || id.Kind == st.FIELD {
					ce.classST.Define(id.Name, id.T, id.Kind)
				} else {
					ce.subroutineST.Define(id.Name, id.T, id.Kind)
				}
			}

			err := ce.CompileDo()
			if err != nil {
				t.Fatalf("CompileDo() error: %v", err)
			}

			vmOutput := vmFile.String()
			want := tt.expectedVMCommands
			if diff := cmp.Diff(want, vmOutput, codeComparer); diff != "" {
				t.Errorf("CompileDo() VM code mismatch (-want +got):\n%s", diff)
				t.Logf("Got VM Code:\n%s", vmOutput)
				t.Logf("Want VM Code:\n%s", want)
			}
		})
	}
}

func TestCompileReturn(t *testing.T) {
	tests := []struct {
		name               string
		jackCode           string
		isVoidFunc         bool
		variables          []st.Identifier // For returning variables
		expectedVMCommands string
	}{
		{
			name:               "Return void",
			jackCode:           `return;`,
			isVoidFunc:         true,
			expectedVMCommands: `push constant 0 return`, // Void functions push 0 before returning
		},
		{
			name:               "Return constant",
			jackCode:           `return 10;`,
			isVoidFunc:         false,
			expectedVMCommands: `push constant 10 return`,
		},
		{
			name:       "Return variable",
			jackCode:   `return i;`,
			isVoidFunc: false,
			variables: []st.Identifier{
				{Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0},
			},
			expectedVMCommands: `push local 0 return`,
		},
		{
			name:       "Return expression",
			jackCode:   `return i + 1;`,
			isVoidFunc: false,
			variables: []st.Identifier{
				{Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0},
			},
			expectedVMCommands: `push local 0 push constant 1 add return`,
		},
		{
			name:               "Return this",
			jackCode:           `return this;`, // Common in constructors
			isVoidFunc:         false,
			expectedVMCommands: `push pointer 0 return`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithVMWriter(vmFile, strings.NewReader(tt.jackCode), "TestClass")

			// Set up context
			funcType := st.NOTVOIDFUNC
			if tt.isVoidFunc {
				funcType = st.VOIDFUNC
			}
			ce.subroutineST.SetCurrentScope("TestClass.testFunc", st.KINDFUNCTION, funcType)

			// Define variables if provided
			for _, id := range tt.variables {
				ce.subroutineST.Define(id.Name, id.T, id.Kind)
			}

			err := ce.CompileReturn()
			if err != nil {
				t.Fatalf("CompileReturn() error: %v", err)
			}

			vmOutput := vmFile.String()
			want := tt.expectedVMCommands
			if diff := cmp.Diff(want, vmOutput, codeComparer); diff != "" {
				t.Errorf("CompileReturn() VM code mismatch (-want +got):\n%s", diff)
				t.Logf("Got VM Code:\n%s", vmOutput)
				t.Logf("Want VM Code:\n%s", want)
			}
		})
	}
}
