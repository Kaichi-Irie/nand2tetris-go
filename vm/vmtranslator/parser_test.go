package vmtranslator

import (
	"bufio"
	"strings"
	"testing"
)

var source string = `
// comment
//// comment
	push   constant 7
// comment

// comment
push constant   8

 add

// comment


////


`

func TestCodeScanner(t *testing.T) {
	cs := CodeScanner{bufio.NewScanner(strings.NewReader(source)), "//"}
	if !cs.scan() {
		t.Errorf("scan() returned false")
	}
	if cs.text() != "push constant 7" {
		t.Errorf("Expected 'push constant 7', got %s", cs.text())
	}
	if !cs.scan() {
		t.Errorf("scan() returned false")
	}
	if cs.text() != "push constant 8" {
		t.Errorf("Expected 'push constant 8', got %s", cs.text())
	}

	if !cs.scan() {
		t.Errorf("scan() returned false")
	}
	if cs.text() != "add" {
		t.Errorf("Expected 'add', got %s", cs.text())
	}
	if cs.scan() {
		t.Errorf("scan() returned true")
	}
}

func TestParser(t *testing.T) {
	cs := CodeScanner{bufio.NewScanner(strings.NewReader(source)), "//"}
	p := Parser{cs, "", C_PUSH}
	if !p.advance() {
		t.Errorf("advance() returned false")
	}
	if p.currentCommand != "push constant 7" {
		t.Errorf("Expected 'push constant 7', got %s", p.currentCommand)
	}
	if p.currentType != C_PUSH {
		t.Errorf("Expected C_PUSH, got %d", p.currentType)
	}

	if !p.advance() {
		t.Errorf("advance() returned false")
	}
	if p.currentCommand != "push constant 8" {
		t.Errorf("Expected 'push constant 8', got %s", p.currentCommand)
	}
	if p.currentType != C_PUSH {
		t.Errorf("Expected C_PUSH, got %d", p.currentType)
	}

	if !p.advance() {
		t.Errorf("advance() returned false")
	}
	if p.currentCommand != "add" {
		t.Errorf("Expected 'add', got %s", p.currentCommand)
	}
	if p.currentType != C_ARITHMETIC {
		t.Errorf("Expected C_ARITHMETIC, got %d", p.currentType)
	}

	if p.advance() {
		t.Errorf("advance() returned true")
	}
}

func TestGetCommandType(t *testing.T) {
	tests := []struct {
		command VMCommand
		want    VMCommandType
	}{
		{"add", C_ARITHMETIC},
		{"sub", C_ARITHMETIC},
		{"neg", C_ARITHMETIC},
		{"eq", C_ARITHMETIC},
		{"gt", C_ARITHMETIC},
		{"lt", C_ARITHMETIC},
		{"and", C_ARITHMETIC},
		{"or", C_ARITHMETIC},
		{"not", C_ARITHMETIC},
		{"push", C_PUSH},
		{"pop", C_POP},
		{"label", C_LABEL},
		{"goto", C_GOTO},
		{"if-goto", C_IF},
		{"function", C_FUNCTION},
		{"return", C_RETURN},
		{"call", C_CALL},
	}
	for _, test := range tests {
		if got := getCommandType(test.command); got != test.want {
			t.Errorf("getCommandType(%s) = %d, want %d", test.command, got, test.want)
		}
	}
}

func TestArg1(t *testing.T) {
	tests := []struct {
		command VMCommand
		want    string
	}{
		{"add", "add"},
		{"sub", "sub"},
		{"neg", "neg"},
		{"eq", "eq"},
		{"gt", "gt"},
		{"lt", "lt"},
		{"and", "and"},
		{"or", "or"},
		{"not", "not"},
		{"push constant 7", "constant"},
		{"pop local 0", "local"},
		{"label LOOP", "LOOP"},
		{"goto LOOP", "LOOP"},
		{"if-goto LOOP", "LOOP"},
		{"function SimpleFunction.test 0", "SimpleFunction.test"},
		{"call SimpleFunction.test 2", "SimpleFunction.test"},
	}
	for _, test := range tests {
		if got := arg1(test.command); got != test.want {
			t.Errorf("arg1(%s) = %s, want %s", test.command, got, test.want)
		}
	}
}

func TestArg2(t *testing.T) {
	tests := []struct {
		command VMCommand
		want    int
	}{
		{"push constant 7", 7},
		{"pop local 0", 0},
		{"function SimpleFunction.test 0", 0},
		{"call SimpleFunction.test 2", 2},
	}
	for _, test := range tests {
		if got := arg2(test.command); got != test.want {
			t.Errorf("arg2(%s) = %d, want %d", test.command, got, test.want)
		}
	}
}
