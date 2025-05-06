package compilationengine

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	st "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/symboltable"
	tk "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/tokenizer"
	vw "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/vmwriter"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	input := strings.NewReader("class Test { }")
	output := &bytes.Buffer{}
	className := "Test"

	engine := New(output, input, className)

	if engine == nil {
		t.Fatal("New() returned nil")
	}
	if engine.vmwriter == nil {
		t.Error("New() did not initialize vmwriter")
	}
	if engine.t == nil {
		t.Error("New() did not initialize tokenizer")
	}
	if engine.classST == nil {
		t.Error("New() did not initialize classST")
	}
	if engine.subroutineST == nil {
		t.Error("New() did not initialize subroutineST")
	}
	if engine.labelCount != 0 {
		t.Errorf("New() labelCount = %d, want 0", engine.labelCount)
	}
}

func TestNewWithFirstToken(t *testing.T) {
	input := strings.NewReader("class Test { }")
	output := &bytes.Buffer{}
	className := "Test"

	engine := NewWithFirstToken(output, input, className)

	if engine == nil {
		t.Fatal("NewWithFirstToken() returned nil")
	}
	if engine.vmwriter == nil {
		t.Error("NewWithFirstToken() did not initialize vmwriter")
	}
	if engine.t == nil {
		t.Fatal("NewWithFirstToken() did not initialize tokenizer")
	}
	if engine.classST == nil {
		t.Error("NewWithFirstToken() did not initialize classST")
	}
	if engine.subroutineST == nil {
		t.Error("NewWithFirstToken() did not initialize subroutineST")
	}
	if engine.labelCount != 0 {
		t.Errorf("NewWithFirstToken() labelCount = %d, want 0", engine.labelCount)
	}

	// Check if the first token was processed correctly
	token := engine.t.CurrentToken
	expectedToken := tk.CLASS
	if token.T != expectedToken.T || token.Val != expectedToken.Val {
		t.Errorf("First token not processed correctly, got %v, want %v", token, expectedToken)
	}
}

func TestNewWithFirstToken_EmptyInput(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	className := "Test"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic with empty input but got none")
		}
	}()

	_ = NewWithFirstToken(output, input, className)
}

// Mock reader that returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestNewWithFirstTokenReaderError(t *testing.T) {
	input := &errorReader{}
	output := &bytes.Buffer{}
	className := "Test"

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic with reader error but got none")
		}
	}()

	_ = NewWithFirstToken(output, input, className)
}

func TestNewWithVMWriter(t *testing.T) {
	input := strings.NewReader("class Test { }")
	output := &bytes.Buffer{}
	className := "Test"

	engine := NewWithVMWriter(output, input, className)

	if engine == nil {
		t.Fatal("NewWithVMWriter() returned nil")
	}
	if engine.vmwriter == nil {
		t.Fatal("NewWithVMWriter() did not initialize vmwriter")
	}
	// Check if vmwriter is correctly initialized (not just the default one from New)
	// This is a bit tricky without exposing internal state or mocking vw.New
	// We assume if it's not nil, it's likely the one created by vw.New in NewWithVMWriter
	if reflect.TypeOf(engine.vmwriter) != reflect.TypeOf(&vw.VMWriter{}) {
		t.Errorf("NewWithVMWriter() did not set vmwriter correctly")
	}

	if engine.t == nil {
		t.Fatal("NewWithVMWriter() did not initialize tokenizer")
	}
	if engine.classST == nil {
		t.Error("NewWithVMWriter() did not initialize classST")
	}
	if engine.subroutineST == nil {
		t.Error("NewWithVMWriter() did not initialize subroutineST")
	}
	if engine.labelCount != 0 {
		t.Errorf("NewWithVMWriter() labelCount = %d, want 0", engine.labelCount)
	}

	// Check if the first token was processed correctly
	token := engine.t.CurrentToken
	expectedToken := tk.CLASS
	if token.T != expectedToken.T || token.Val != expectedToken.Val {
		t.Errorf("First token not processed correctly, got %v, want %v", token, expectedToken)
	}
}

func TestCompilationEngineStructure(t *testing.T) {
	engineType := reflect.TypeOf(CompilationEngine{})

	fields := map[string]reflect.Type{
		"vmwriter":     reflect.TypeOf(&vw.VMWriter{}),
		"t":            reflect.TypeOf(&tk.Tokenizer{}),
		"classST":      reflect.TypeOf(&st.SymbolTable{}),
		"subroutineST": reflect.TypeOf(&st.SymbolTable{}),
		"labelCount":   reflect.TypeOf(int(0)),
	}

	for name, expectedType := range fields {
		field, exists := engineType.FieldByName(name)
		if !exists {
			t.Errorf("CompilationEngine missing '%s' field", name)
			continue
		}
		if field.Type != expectedType {
			t.Errorf("'%s' field has incorrect type: got %v, want %v", name, field.Type, expectedType)
		}
	}
}

func TestLookup(t *testing.T) {
	input := strings.NewReader("")
	output := &bytes.Buffer{}
	className := "Test"
	engine := New(output, input, className)

	// Setup symbol tables
	classVar := st.Identifier{Name: "classVar", Kind: st.STATIC, T: "int", Index: 0}
	subroutineVar := st.Identifier{Name: "subVar", Kind: st.VAR, T: "char", Index: 0}
	argVar := st.Identifier{Name: "argVar", Kind: st.ARG, T: "boolean", Index: 0}

	engine.classST.Define(classVar.Name, classVar.T, classVar.Kind)
	engine.subroutineST.Define(subroutineVar.Name, subroutineVar.T, subroutineVar.Kind)
	engine.subroutineST.Define(argVar.Name, argVar.T, argVar.Kind)

	tests := []struct {
		name         string
		wantID       st.Identifier
		wantOK       bool
		setupSubOnly bool // Flag to reset classST for specific tests
	}{
		{"subVar", subroutineVar, true, false},
		{"argVar", argVar, true, false},
		{"classVar", classVar, true, false},
		{"nonExistent", st.Identifier{}, false, false},
		{"classVarSubOnly", st.Identifier{}, false, true}, // Test when classVar is not found because classST is reset
	}

	originalClassST := *engine.classST // Save original classST

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupSubOnly {
				engine.classST = st.NewSymbolTable() // Reset classST for this test case
			} else {
				engine.classST = &originalClassST // Restore original classST
			}

			gotID, gotOK := engine.Lookup(tt.name)

			if gotOK != tt.wantOK {
				t.Errorf("Lookup(%q) ok = %v, want %v", tt.name, gotOK, tt.wantOK)
			}
			// Use cmp.Equal for struct comparison
			if !cmp.Equal(gotID, tt.wantID) {
				t.Errorf("Lookup(%q) id diff = %v", tt.name, cmp.Diff(tt.wantID, gotID))
			}
		})
	}
	engine.classST = &originalClassST // Ensure classST is restored after all tests
}
