package compilationengine

import (
	"bytes"
	"io"
	tk "nand2tetris-go/jackcompiler/tokenizer"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	// Test with valid reader and writer
	input := strings.NewReader("class Test { }")
	output := &bytes.Buffer{}

	engine := New(output, input)

	// Verify engine is properly initialized
	if engine == nil {
		t.Fatal("New() returned nil")
	}

	if engine.writer != output {
		t.Errorf("writer not properly set, got %v, want %v", engine.writer, output)
	}

	// Verify tokenizer was created
	if engine.t == nil {
		t.Fatal("tokenizer was not initialized")
	}
}

func TestNewWithFirstToken(t *testing.T) {
	// Test with valid input that contains a token
	input := strings.NewReader("class Test { }")
	output := &bytes.Buffer{}

	engine := NewWithFirstToken(output, input)

	// Verify engine is properly initialized
	if engine == nil {
		t.Fatal("NewWithFirstToken() returned nil")
	}

	if engine.writer != output {
		t.Errorf("writer not properly set, got %v, want %v", engine.writer, output)
	}

	// Verify tokenizer was created and has the first token ready
	if engine.t == nil {
		t.Fatal("tokenizer was not initialized")
	}

	// Check if the first token was processed correctly
	token := engine.t.CurrentToken
	if token.T != tk.TT_KEYWORD || token.Val != tk.CLASS.Val {
		t.Errorf("First token not processed correctly, got %v, want 'class'", token)
	}
}

func TestNewWithFirstToken_EmptyInput(t *testing.T) {
	// Test behavior with empty input
	input := strings.NewReader("")
	output := &bytes.Buffer{}

	// Should panic with empty input - use recover to test this
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic with empty input but got none")
		}
	}()

	_ = NewWithFirstToken(output, input)
}

// Mock reader that returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestNewWithFirstToken_ReaderError(t *testing.T) {
	// Test behavior with a reader that returns errors
	input := &errorReader{}
	output := &bytes.Buffer{}

	// Should panic with reader error - use recover to test this
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic with reader error but got none")
		}
	}()

	_ = NewWithFirstToken(output, input)
}

// Test to ensure the CompilationEngine struct has expected fields and types
func TestCompilationEngineStructure(t *testing.T) {
	// Use reflection to inspect the struct
	engineType := reflect.TypeOf(CompilationEngine{})

	// Check writer field
	writerField, exists := engineType.FieldByName("writer")
	if !exists {
		t.Fatal("CompilationEngine missing 'writer' field")
	}
	if writerField.Type != reflect.TypeOf((*io.Writer)(nil)).Elem() {
		t.Errorf("writer field has incorrect type: got %v, want io.Writer", writerField.Type)
	}

	// Check tokenizer field
	tokenizerField, exists := engineType.FieldByName("t")
	if !exists {
		t.Fatal("CompilationEngine missing 't' field")
	}
	if tokenizerField.Type != reflect.TypeOf((*tk.Tokenizer)(nil)) {
		t.Errorf("t field has incorrect type: got %v, want *tokenizer.Tokenizer", tokenizerField.Type)
	}
}
