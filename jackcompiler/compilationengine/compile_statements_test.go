package compilationengine

import (
	"bytes"
	st "nand2tetris-go/jackcompiler/symboltable"
	tk "nand2tetris-go/jackcompiler/tokenizer"
	"strings"
	"testing"

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
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode), "")
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

func TestCompileLet2(t *testing.T) {
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
	}

	for _, test := range tests {
		vmFile := &bytes.Buffer{}
		ce := NewWithVMWriter(vmFile, &bytes.Buffer{}, strings.NewReader(test.jackCode), "")

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
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode), "")
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
func TestCompileIf2(t *testing.T) {
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
		ce := NewWithVMWriter(vmFile, &bytes.Buffer{}, strings.NewReader(test.jackCode), "")

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
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode), "")
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
func TestCompileWhile2(t *testing.T) {
	variables := []st.Identifier{
		{Name: "i", Kind: st.VAR, T: tk.INT.Val, Index: 0}}
	tests := []struct {
		jackCode          string
		expectedVMCommand string
	}{
		{
			jackCode: `while (i) { }`,
			expectedVMCommand: `label label0
			push local 0
			not
			if-goto label1
			goto label0
			label label1`},
		{
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
	for _, test := range tests {
		vmFile := &bytes.Buffer{}
		ce := NewWithVMWriter(vmFile, &bytes.Buffer{}, strings.NewReader(test.jackCode), "")

		// Define the variables in the symbol table
		for _, id := range variables {
			ce.classST.Define(id.Name, id.T, id.Kind)
		}
		err := ce.CompileWhile()
		if err != nil {
			t.Errorf("CompileWhile() error: %v", err)
		}

		// trim leading and trailing whitespace
		vmOutput := vmFile.String()
		want := test.expectedVMCommand
		if diff := cmp.Diff(vmOutput, want, codeComparer); diff != "" {
			t.Errorf("CompileWhile() mismatch (-got +want):\n%s\n%s", vmOutput, want)
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
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode), "")
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
		ce := NewWithFirstToken(xmlFile, strings.NewReader(test.jackCode), "")
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
