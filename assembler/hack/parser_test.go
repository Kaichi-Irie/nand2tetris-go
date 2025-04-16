package hack

import (
	"bufio"
	"strings"
	"testing"
)

func TestParserAdvance(t *testing.T) {
	input := `
// This is a comment
@2
D=A
(LOOP)
0;JMP
`
	scanner := bufio.NewScanner(strings.NewReader(input))
	parser := &Parser{scanner: scanner}

	tests := []struct {
		expectedInstruction Instruction
		expectedType        InstructionType
	}{
		{"@2", A_Instruction},
		{"D=A", C_Instruction},
		{"(LOOP)", L_Instruction},
		{"0;JMP", C_Instruction},
	}

	for _, tt := range tests {
		if !parser.advance() {
			t.Fatalf("Expected more instructions, but advance returned false")
		}
		if parser.currentInstruction != tt.expectedInstruction {
			t.Errorf("Expected instruction %q, got %q", tt.expectedInstruction, parser.currentInstruction)
		}
		if parser.currentType != tt.expectedType {
			t.Errorf("Expected type %v, got %v", tt.expectedType, parser.currentType)
		}
	}

	if parser.advance() {
		t.Errorf("Expected no more instructions, but advance returned true")
	}
}

func TestParserSymbol(t *testing.T) {
	input := `@100
(LOOP)
D=A
`
	scanner := bufio.NewScanner(strings.NewReader(input))
	parser := &Parser{scanner: scanner}

	tests := []struct {
		expectedSymbol SymbolOrConstant
		shouldError    bool
	}{
		{"100", false},
		{"LOOP", false},
		{"", true}, // C instruction should return an error
	}

	for _, tt := range tests {
		parser.advance()
		symbol, err := parser.symbol()
		if tt.shouldError {
			if err == nil {
				t.Errorf("Expected an error, but got none")
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if symbol != tt.expectedSymbol {
				t.Errorf("Expected symbol %q, got %q", tt.expectedSymbol, symbol)
			}
		}
	}
}

func TestParserDest(t *testing.T) {
	input := `D=A
0;JMP
@2
`
	scanner := bufio.NewScanner(strings.NewReader(input))
	parser := &Parser{scanner: scanner}

	tests := []struct {
		expectedDest Mnemonic
		shouldError  bool
	}{
		{"D", false},
		{"null", false},
		{"", true}, // A instruction should return an error
	}

	for _, tt := range tests {
		parser.advance()
		dest, err := parser.dest()
		if tt.shouldError {
			if err == nil {
				t.Errorf("Expected an error, but got none")
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if dest != tt.expectedDest {
				t.Errorf("Expected dest %q, got %q", tt.expectedDest, dest)
			}
		}
	}
}

func TestParserComp(t *testing.T) {
	input := `D=A
0;JMP
@2
`
	scanner := bufio.NewScanner(strings.NewReader(input))
	parser := &Parser{scanner: scanner}

	tests := []struct {
		expectedComp Mnemonic
		shouldError  bool
	}{
		{"A", false},
		{"0", false},
		{"", true}, // A instruction should return an error
	}

	for _, tt := range tests {
		parser.advance()
		comp, err := parser.comp()
		if tt.shouldError {
			if err == nil {
				t.Errorf("Expected an error, but got none")
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if comp != tt.expectedComp {
				t.Errorf("Expected comp %q, got %q", tt.expectedComp, comp)
			}
		}
	}
}

func TestParserJump(t *testing.T) {
	input := `D=A
0;JMP
@2
`
	scanner := bufio.NewScanner(strings.NewReader(input))
	parser := &Parser{scanner: scanner}

	tests := []struct {
		expectedJump Mnemonic
		shouldError  bool
	}{
		{"null", false},
		{"JMP", false},
		{"", true}, // A instruction should return an error
	}

	for _, tt := range tests {
		parser.advance()
		jump, err := parser.jump()
		if tt.shouldError {
			if err == nil {
				t.Errorf("Expected an error, but got none")
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if jump != tt.expectedJump {
				t.Errorf("Expected jump %q, got %q", tt.expectedJump, jump)
			}
		}
	}
}
