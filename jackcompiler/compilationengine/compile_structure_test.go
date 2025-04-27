package compilationengine

import (
	"bytes"
	st "nand2tetris-go/jackcompiler/symboltable"
	tk "nand2tetris-go/jackcompiler/tokenizer"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// codeComparer defined in compile_statements_test.go is assumed to be available
// If not, add it here:
/*
var codeComparer = cmp.Comparer(func(x string, y string) bool {
    x = strings.TrimSpace(x)
    y = strings.TrimSpace(y)
    x = strings.Join(strings.Fields(x), " ")
    y = strings.Join(strings.Fields(y), " ")
    return x == y
})
*/

func TestCompileClass(t *testing.T) {
	tests := []struct {
		name                 string
		jackCode             string
		className            string
		expectedClassST      st.SymbolTable
		expectedSubroutineST st.SymbolTable // Expected state *after* the last subroutine if any
		expectedVMCode       string
	}{
		{
			name:      "Class with static vars only",
			jackCode:  `class Main { static int i; static int j; }`,
			className: "Main",
			expectedClassST: st.SymbolTable{
				CurrentScope: st.Identifier{Name: "Main", Kind: st.KINDCLASS, T: st.NOTVOIDFUNC, Index: -1},
				StaticCnt:    2,
				VariableMap: map[string]st.Identifier{
					"Main": {Name: "Main", Kind: st.NONE, T: "Main", Index: 0},
					"i":    {Name: "i", Kind: st.STATIC, T: tk.INT.Val, Index: 0},
					"j":    {Name: "j", Kind: st.STATIC, T: tk.INT.Val, Index: 1},
				},
			},
			expectedVMCode: ``, // No VM code generated directly by class/var decs
		},
		{
			name:      "Class with field vars only",
			jackCode:  `class Point { field int x, y; }`,
			className: "Point",
			expectedClassST: st.SymbolTable{
				CurrentScope: st.Identifier{Name: "Point", Kind: st.KINDCLASS, T: st.NOTVOIDFUNC, Index: -1},
				FieldCnt:     2,
				VariableMap: map[string]st.Identifier{
					"Point": {Name: "Point", Kind: st.NONE, T: "Point", Index: 0},
					"x":     {Name: "x", Kind: st.FIELD, T: tk.INT.Val, Index: 0},
					"y":     {Name: "y", Kind: st.FIELD, T: tk.INT.Val, Index: 1},
				},
			},
			expectedVMCode: ``, // No VM code generated directly by class/var decs
		},
		{
			name: "Class with fields and a simple function",
			jackCode: `class Main {
                field int i, j;
                function void main() {
                    return;
                }
            }`,
			className: "Main",
			expectedClassST: st.SymbolTable{
				CurrentScope: st.Identifier{Name: "Main", Kind: st.KINDCLASS, T: st.NOTVOIDFUNC, Index: -1},
				FieldCnt:     2,
				VariableMap: map[string]st.Identifier{
					"Main": {Name: "Main", Kind: st.NONE, T: "Main", Index: 0},
					"i":    {Name: "i", Kind: st.FIELD, T: tk.INT.Val, Index: 0},
					"j":    {Name: "j", Kind: st.FIELD, T: tk.INT.Val, Index: 1},
				},
			},
			expectedSubroutineST: st.SymbolTable{ // State after main finishes
				CurrentScope: st.Identifier{Name: "Main.main", Kind: st.KINDFUNCTION, T: st.VOIDFUNC, Index: -1},
				VariableMap:  map[string]st.Identifier{},
			},
			expectedVMCode: `function Main.main 0
push constant 0
return
`,
		},
		{
			name: "Class with fields, function with args and vars",
			jackCode: `class Main {
                field MyClass i, j;
                function void main(int a, boolean b) {
                    var int c;
                    var boolean d, e;
                    return;
                }
            }`,
			className: "Main",
			expectedClassST: st.SymbolTable{
				CurrentScope: st.Identifier{Name: "Main", Kind: st.KINDCLASS, T: st.NOTVOIDFUNC, Index: -1},
				FieldCnt:     2,
				VariableMap: map[string]st.Identifier{
					"Main": {Name: "Main", Kind: st.NONE, T: "Main", Index: 0},
					"i":    {Name: "i", Kind: st.FIELD, T: "MyClass", Index: 0},
					"j":    {Name: "j", Kind: st.FIELD, T: "MyClass", Index: 1},
				},
			},
			expectedSubroutineST: st.SymbolTable{ // State after main finishes
				CurrentScope: st.Identifier{Name: "Main.main", Kind: st.KINDFUNCTION, T: st.VOIDFUNC, Index: -1},
				ArgCnt:       2,
				VarCnt:       3,
				VariableMap: map[string]st.Identifier{
					"a": {Name: "a", Kind: st.ARG, T: "int", Index: 0},
					"b": {Name: "b", Kind: st.ARG, T: "boolean", Index: 1},
					"c": {Name: "c", Kind: st.VAR, T: "int", Index: 0},
					"d": {Name: "d", Kind: st.VAR, T: "boolean", Index: 1},
					"e": {Name: "e", Kind: st.VAR, T: "boolean", Index: 2},
				},
			},
			expectedVMCode: `function Main.main 3
push constant 0
return
`,
		},
		{
			name: "Class with constructor",
			jackCode: `class Point {
                field int x, y;
                constructor Point new(int ax, int ay) {
                    let x = ax;
                    let y = ay;
                    return this;
                }
            }`,
			className: "Point",
			expectedClassST: st.SymbolTable{
				CurrentScope: st.Identifier{Name: "Point", Kind: st.KINDCLASS, T: st.NOTVOIDFUNC, Index: -1},
				FieldCnt:     2,
				VariableMap: map[string]st.Identifier{
					"Point": {Name: "Point", Kind: st.NONE, T: "Point", Index: 0},
					"x":     {Name: "x", Kind: st.FIELD, T: tk.INT.Val, Index: 0},
					"y":     {Name: "y", Kind: st.FIELD, T: tk.INT.Val, Index: 1},
				},
			},
			expectedSubroutineST: st.SymbolTable{ // State after constructor finishes
				CurrentScope: st.Identifier{Name: "Point.new", Kind: st.KINDCONSTRUCTOR, T: st.NOTVOIDFUNC, Index: -1},
				ArgCnt:       2,
				VariableMap: map[string]st.Identifier{
					"ax": {Name: "ax", Kind: st.ARG, T: "int", Index: 0},
					"ay": {Name: "ay", Kind: st.ARG, T: "int", Index: 1},
				},
			},
			expectedVMCode: `function Point.new 0
push constant 2
call Memory.alloc 1
pop pointer 0
push argument 0
pop this 0
push argument 1
pop this 1
push pointer 0
return
`,
		},
		{
			name: "Class with method",
			jackCode: `class Point {
                field int x, y;
                method int getX() {
                    return x;
                }
            }`,
			className: "Point",
			expectedClassST: st.SymbolTable{
				CurrentScope: st.Identifier{Name: "Point", Kind: st.KINDCLASS, T: st.NOTVOIDFUNC, Index: -1},
				FieldCnt:     2,
				VariableMap: map[string]st.Identifier{
					"Point": {Name: "Point", Kind: st.NONE, T: "Point", Index: 0},
					"x":     {Name: "x", Kind: st.FIELD, T: tk.INT.Val, Index: 0},
					"y":     {Name: "y", Kind: st.FIELD, T: tk.INT.Val, Index: 1},
				},
			},
			expectedSubroutineST: st.SymbolTable{ // State after method finishes
				CurrentScope: st.Identifier{Name: "Point.getX", Kind: st.KINDMETHOD, T: st.NOTVOIDFUNC, Index: -1},
				ArgCnt:       1, // Includes 'this'
				VariableMap: map[string]st.Identifier{
					"this": {Name: "this", Kind: st.ARG, T: "Point", Index: 0},
				},
			},
			expectedVMCode: `function Point.getX 0
push argument 0
pop pointer 0
push this 0
return
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			// Use NewWithVMWriter to capture VM output
			ce := NewWithVMWriter(vmFile, strings.NewReader(tt.jackCode), tt.className)
			err := ce.CompileClass()
			if err != nil {
				t.Fatalf("CompileClass() error: %v", err)
			}

			// Compare Class Symbol Table
			if diff := cmp.Diff(tt.expectedClassST, *ce.classST); diff != "" {
				t.Errorf("Class Symbol Table mismatch (-want +got):\n%s", diff)
			}

			// Compare Subroutine Symbol Table (if applicable)
			// Note: This checks the state *after* the last subroutine compiled
			if !cmp.Equal(tt.expectedSubroutineST, st.SymbolTable{}) { // Only check if expected is not empty
				if diff := cmp.Diff(tt.expectedSubroutineST, *ce.subroutineST); diff != "" {
					t.Errorf("Subroutine Symbol Table mismatch (-want +got):\n%s", diff)
				}
			}

			// Compare VM Code
			gotVMCode := vmFile.String()
			if diff := cmp.Diff(tt.expectedVMCode, gotVMCode, codeComparer); diff != "" {
				t.Errorf("VM Code mismatch (-want +got):\n%s", diff)
				// For easier debugging:
				t.Logf("Got VM Code:\n%s", gotVMCode)
				t.Logf("Want VM Code:\n%s", tt.expectedVMCode)
			}
		})
	}
}

// TestCompileClassVarDec focuses *only* on the symbol table state,
// as this function itself doesn't directly generate VM code.
func TestCompileClassVarDec(t *testing.T) {
	tests := []struct {
		name                string
		jackCode            string // Code snippet starting with static/field
		fieldOrStatic       tk.Token
		expectedSymbolTable st.SymbolTable
	}{
		{
			name:          "Single static int",
			jackCode:      `static int i;`,
			fieldOrStatic: tk.STATIC,
			expectedSymbolTable: st.SymbolTable{
				StaticCnt: 1,
				VariableMap: map[string]st.Identifier{
					"i": {Name: "i", Kind: st.STATIC, T: tk.INT.Val, Index: 0},
				},
			},
		},
		{
			name:          "Multiple field int",
			jackCode:      `field int i, j;`,
			fieldOrStatic: tk.FIELD,
			expectedSymbolTable: st.SymbolTable{
				FieldCnt: 2,
				VariableMap: map[string]st.Identifier{
					"i": {Name: "i", Kind: st.FIELD, T: tk.INT.Val, Index: 0},
					"j": {Name: "j", Kind: st.FIELD, T: tk.INT.Val, Index: 1},
				},
			},
		},
		{
			name:          "Multiple static custom type",
			jackCode:      `static MyClass a, b;`,
			fieldOrStatic: tk.STATIC,
			expectedSymbolTable: st.SymbolTable{
				StaticCnt: 2,
				VariableMap: map[string]st.Identifier{
					"a": {Name: "a", Kind: st.STATIC, T: "MyClass", Index: 0},
					"b": {Name: "b", Kind: st.STATIC, T: "MyClass", Index: 1},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{} // Not used for checking output here
			// Need to use NewWithFirstToken because CompileClassVarDec expects the token to be ready
			ce := NewWithFirstToken(vmFile, strings.NewReader(tt.jackCode), "TestClass")

			// Manually set the class scope for context
			ce.classST.SetCurrentScope("TestClass", st.KINDCLASS, st.NOTVOIDFUNC)

			err := ce.CompileClassVarDec(tt.fieldOrStatic)
			if err != nil {
				t.Fatalf("CompileClassVarDec() error: %v", err)
			}

			// Compare only the relevant parts of the symbol table
			gotST := *ce.classST
			wantST := tt.expectedSymbolTable

			if gotST.StaticCnt != wantST.StaticCnt {
				t.Errorf("StaticCnt mismatch: got %d, want %d", gotST.StaticCnt, wantST.StaticCnt)
			}
			if gotST.FieldCnt != wantST.FieldCnt {
				t.Errorf("FieldCnt mismatch: got %d, want %d", gotST.FieldCnt, wantST.FieldCnt)
			}
			if diff := cmp.Diff(wantST.VariableMap, gotST.VariableMap); diff != "" {
				t.Errorf("VariableMap mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestCompileSubroutine focuses on VM code and subroutine symbol table.
func TestCompileSubroutine(t *testing.T) {
	tests := []struct {
		name                 string
		jackCode             string // Code snippet for the subroutine
		className            string
		expectedSubroutineST st.SymbolTable
		expectedVMCode       string
	}{
		{
			name:      "Simple void function",
			jackCode:  `function void main() { var int i; return; }`,
			className: "Test",
			expectedSubroutineST: st.SymbolTable{
				CurrentScope: st.Identifier{Name: "Test.main", Kind: st.KINDFUNCTION, T: st.VOIDFUNC, Index: -1},
				VarCnt:       1,
				VariableMap: map[string]st.Identifier{
					"i": {Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0},
				},
			},
			expectedVMCode: `function Test.main 1
push constant 0
return
`,
		},
		{
			name:      "Constructor with params",
			jackCode:  `constructor Point new(int ax) { var int x; let x = ax; return this; }`,
			className: "Point",
			expectedSubroutineST: st.SymbolTable{
				CurrentScope: st.Identifier{Name: "Point.new", Kind: st.KINDCONSTRUCTOR, T: st.NOTVOIDFUNC, Index: -1},
				ArgCnt:       1,
				VarCnt:       1,
				VariableMap: map[string]st.Identifier{
					"ax": {Name: "ax", Kind: st.ARG, T: tk.INT.Val, Index: 0},
					"x":  {Name: "x", Kind: st.VAR, T: tk.INT.Val, Index: 0},
				},
			},
			// Assumes classST has FieldCnt=1 for 'x' before calling CompileSubroutine
			expectedVMCode: `function Point.new 1
push constant 1
call Memory.alloc 1
pop pointer 0
push argument 0
pop local 0
push pointer 0
return
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			ce := NewWithVMWriter(vmFile, strings.NewReader(tt.jackCode), tt.className)

			// Set up class context (needed for method/constructor this/alloc)
			ce.classST.SetCurrentScope(tt.className, st.KINDCLASS, st.NOTVOIDFUNC)
			// Manually set FieldCnt if needed for constructor test
			if strings.Contains(tt.jackCode, "constructor") {
				// Crude way to estimate FieldCnt for the test case
				ce.classST.FieldCnt = 1 // Adjust if test case implies more fields
			}
			// Manually define class fields if needed for method test
			if strings.Contains(tt.jackCode, "method") {
				ce.classST.Define("x", "int", st.FIELD) // Adjust if test case implies different fields
			}

			err := ce.CompileSubroutine()
			if err != nil {
				t.Fatalf("CompileSubroutine() error: %v", err)
			}

			// Compare Subroutine Symbol Table
			if diff := cmp.Diff(tt.expectedSubroutineST, *ce.subroutineST); diff != "" {
				t.Errorf("Subroutine Symbol Table mismatch (-want +got):\n%s", diff)
			}

			// Compare VM Code
			gotVMCode := vmFile.String()
			if diff := cmp.Diff(tt.expectedVMCode, gotVMCode, codeComparer); diff != "" {
				t.Errorf("VM Code mismatch (-want +got):\n%s", diff)
				t.Logf("Got VM Code:\n%s", gotVMCode)
				t.Logf("Want VM Code:\n%s", tt.expectedVMCode)
			}
		})
	}
}

// TestCompileVarDec focuses *only* on the symbol table state.
func TestCompileVarDec(t *testing.T) {
	tests := []struct {
		name                string
		jackCode            string // Code snippet starting with 'var'
		expectedSymbolTable st.SymbolTable
	}{
		{
			name:     "Single var int",
			jackCode: `var int i;`,
			expectedSymbolTable: st.SymbolTable{
				VarCnt: 1,
				VariableMap: map[string]st.Identifier{
					"i": {Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0},
				},
			},
		},
		{
			name:     "Multiple var boolean",
			jackCode: `var boolean flag, done;`,
			expectedSymbolTable: st.SymbolTable{
				VarCnt: 2,
				VariableMap: map[string]st.Identifier{
					"flag": {Name: "flag", Kind: st.VAR, T: tk.BOOLEAN.Val, Index: 0},
					"done": {Name: "done", Kind: st.VAR, T: tk.BOOLEAN.Val, Index: 1},
				},
			},
		},
		{
			name:     "Var custom type",
			jackCode: `var MyObject obj;`,
			expectedSymbolTable: st.SymbolTable{
				VarCnt: 1,
				VariableMap: map[string]st.Identifier{
					"obj": {Name: "obj", Kind: st.VAR, T: "MyObject", Index: 0},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{} // Not used
			ce := NewWithFirstToken(vmFile, strings.NewReader(tt.jackCode), "TestClass")

			// Set up subroutine context
			ce.subroutineST.SetCurrentScope("Test.testFunc", st.KINDFUNCTION, st.VOIDFUNC)

			err := ce.CompileVarDec()
			if err != nil {
				t.Fatalf("CompileVarDec() error: %v", err)
			}

			// Compare relevant parts of the symbol table
			gotST := *ce.subroutineST
			wantST := tt.expectedSymbolTable

			if gotST.VarCnt != wantST.VarCnt {
				t.Errorf("VarCnt mismatch: got %d, want %d", gotST.VarCnt, wantST.VarCnt)
			}
			if diff := cmp.Diff(wantST.VariableMap, gotST.VariableMap); diff != "" {
				t.Errorf("VariableMap mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestCompileParameterList focuses *only* on the symbol table state.
func TestCompileParameterList(t *testing.T) {
	tests := []struct {
		name                string
		jackCode            string // Code snippet for the parameter list content (inside parens)
		isMethod            bool   // To check if 'this' is added correctly
		className           string // Needed if isMethod is true
		expectedSymbolTable st.SymbolTable
	}{
		{
			name:     "Empty list",
			jackCode: ``,
			expectedSymbolTable: st.SymbolTable{
				ArgCnt:      0,
				VariableMap: map[string]st.Identifier{},
			},
		},
		{
			name:     "Single int param",
			jackCode: `int i`,
			expectedSymbolTable: st.SymbolTable{
				ArgCnt: 1,
				VariableMap: map[string]st.Identifier{
					"i": {Name: "i", Kind: st.ARG, T: tk.INT.Val, Index: 0},
				},
			},
		},
		{
			name:     "Multiple mixed params",
			jackCode: `int i, boolean b, MyClass c`,
			expectedSymbolTable: st.SymbolTable{
				ArgCnt: 3,
				VariableMap: map[string]st.Identifier{
					"i": {Name: "i", Kind: st.ARG, T: tk.INT.Val, Index: 0},
					"b": {Name: "b", Kind: st.ARG, T: tk.BOOLEAN.Val, Index: 1},
					"c": {Name: "c", Kind: st.ARG, T: "MyClass", Index: 2},
				},
			},
		},
		{
			name:      "Method empty list (adds this)",
			jackCode:  ``,
			isMethod:  true,
			className: "MyObj",
			expectedSymbolTable: st.SymbolTable{
				ArgCnt: 1,
				VariableMap: map[string]st.Identifier{
					"this": {Name: "this", Kind: st.ARG, T: "MyObj", Index: 0},
				},
			},
		},
		{
			name:      "Method with params (adds this)",
			jackCode:  `int count`,
			isMethod:  true,
			className: "Counter",
			expectedSymbolTable: st.SymbolTable{
				ArgCnt: 2,
				VariableMap: map[string]st.Identifier{
					"this":  {Name: "this", Kind: st.ARG, T: "Counter", Index: 0},
					"count": {Name: "count", Kind: st.ARG, T: tk.INT.Val, Index: 1},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{} // Not used
			// Need to wrap the code slightly for the tokenizer
			reader := strings.NewReader(tt.jackCode + ")") // Add closing paren for context
			ce := NewWithFirstToken(vmFile, reader, "TestClass")

			// Set up subroutine context
			kind := st.KINDFUNCTION
			if tt.isMethod {
				kind = st.KINDMETHOD
			}
			ce.subroutineST.SetCurrentScope("Test.testFunc", kind, st.VOIDFUNC)

			// Manually add 'this' if it's a method *before* calling CompileParameterList
			if tt.isMethod {
				ce.subroutineST.Define("this", tt.className, st.ARG)
			}

			err := ce.CompileParameterList()
			if err != nil {
				t.Fatalf("CompileParameterList() error: %v", err)
			}

			// Compare relevant parts of the symbol table
			gotST := *ce.subroutineST
			wantST := tt.expectedSymbolTable

			if gotST.ArgCnt != wantST.ArgCnt {
				t.Errorf("ArgCnt mismatch: got %d, want %d", gotST.ArgCnt, wantST.ArgCnt)
			}
			if diff := cmp.Diff(wantST.VariableMap, gotST.VariableMap); diff != "" {
				t.Errorf("VariableMap mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestCompileSubroutineBody focuses on VM code generation within the body.
func TestCompileSubroutineBody(t *testing.T) {
	tests := []struct {
		name           string
		jackCode       string // Code snippet for the body content (inside braces)
		subroutineKind string // function, method, constructor
		className      string // Needed for method/constructor
		fieldCount     int    // Needed for constructor alloc
		varCount       int    // Expected var count for function call
		expectedVMCode string
	}{
		{
			name:           "Empty body (function)",
			jackCode:       `{ }`,
			subroutineKind: st.KINDFUNCTION,
			className:      "Test",
			varCount:       0,
			expectedVMCode: `function Test.testFunc 0
`, // Function call is written here
		},
		{
			name:           "Body with vars (function)",
			jackCode:       `{ var int i; var char c; }`,
			subroutineKind: st.KINDFUNCTION,
			className:      "Test",
			varCount:       2,
			expectedVMCode: `function Test.testFunc 2
`, // Statements would follow
		},
		{
			name:           "Body with statements (function)",
			jackCode:       `{ return; }`,
			subroutineKind: st.KINDFUNCTION,
			className:      "Test",
			varCount:       0,
			expectedVMCode: `function Test.testFunc 0
push constant 0
return
`,
		},

		{
			name:           "Body with statements (constructor)",
			jackCode:       `{ return this; }`, // Assumes 1 field
			subroutineKind: st.KINDCONSTRUCTOR,
			className:      "Thing",
			fieldCount:     1,
			varCount:       0,
			expectedVMCode: `function Thing.new 0
push constant 1
call Memory.alloc 1
pop pointer 0
push pointer 0
return
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vmFile := &bytes.Buffer{}
			// Need to wrap code slightly for tokenizer
			reader := strings.NewReader(tt.jackCode)
			ce := NewWithVMWriter(vmFile, reader, tt.className)

			// Set up context
			ce.classST.SetCurrentScope(tt.className, st.KINDCLASS, st.NOTVOIDFUNC)
			ce.classST.FieldCnt = tt.fieldCount
			// Define fields if method test assumes them
			if tt.subroutineKind == st.KINDMETHOD {
				ce.classST.Define("x", "int", st.FIELD)
				ce.classST.Define("y", "int", st.FIELD)
			}

			subroutineName := tt.className + "."
			if tt.subroutineKind == st.KINDCONSTRUCTOR {
				subroutineName += "new"
			} else if tt.subroutineKind == st.KINDMETHOD {
				subroutineName += "testMethod" // Generic name for testing
				if tt.name == "Complex body (method)" {
					subroutineName = "Point.complexMethod"
				}
			} else {
				subroutineName += "testFunc" // Generic name for testing
			}

			ce.subroutineST.SetCurrentScope(subroutineName, tt.subroutineKind, st.VOIDFUNC) // Assume void unless return value implies otherwise
			// Manually define vars if needed *before* calling CompileSubroutineBody
			if strings.Contains(tt.jackCode, "var int temp") {
				ce.subroutineST.Define("temp", "int", st.VAR)
			}
			if tt.subroutineKind == st.KINDMETHOD {
				ce.subroutineST.Define("this", tt.className, st.ARG) // Methods have 'this'
			}

			err := ce.CompileSubroutineBody()
			if err != nil {
				t.Fatalf("CompileSubroutineBody() error: %v", err)
			}

			// Compare VM Code
			gotVMCode := vmFile.String()
			if diff := cmp.Diff(tt.expectedVMCode, gotVMCode, codeComparer); diff != "" {
				t.Errorf("VM Code mismatch (-want +got):\n%s", diff)
				t.Logf("Got VM Code:\n%s", gotVMCode)
				t.Logf("Want VM Code:\n%s", tt.expectedVMCode)
			}
		})
	}
}
