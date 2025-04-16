package hack

import (
	"testing"
)

func TestDecimalToBinary(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"0", "000000000000000", false},
		{"1", "000000000000001", false},
		{"100", "000000001100100", false},
		{"32767", "111111111111111", false}, // Max 15-bit value
		{"abc", "", true},                   // Negative number
	}

	for _, test := range tests {
		code, err := decimalToBinary(test.input)
		if (err != nil) != test.hasError || string(code) != test.expected {
			t.Errorf(`decimalToBinary("%s") = %q, %v, want %q, error: %v`, test.input, code, err, test.expected, test.hasError)
		}
	}
}

func TestSymbol(t *testing.T) {
	st := NewSymbolTable()
	st.table["LABEL"] = 8

	tests := []struct {
		input    SymbolOrConstant
		expected string
		hasError bool
	}{
		{"100", "000000001100100", false},     // Constant
		{"LABEL", "000000000001000", false},   // Existing symbol
		{"NEW_VAR", "000000000010000", false}, // New variable
	}

	for _, test := range tests {
		code, err := symbol(test.input, &st)
		if (err != nil) != test.hasError || string(code) != test.expected {
			t.Errorf(`symbol("%s") = %q, %v, want %q, error: %v`, test.input, code, err, test.expected, test.hasError)
		}
	}
}

func TestDest(t *testing.T) {
	tests := []struct {
		input    Mnemonic
		expected string
		hasError bool
	}{
		{"null", "000", false},
		{"M", "001", false},
		{"D", "010", false},
		{"MD", "011", false},
		{"A", "100", false},
		{"AM", "101", false},
		{"AD", "110", false},
		{"AMD", "111", false},
		{"INVALID", "", true}, // Invalid mnemonic
	}

	for _, test := range tests {
		code, err := dest(test.input)
		if (err != nil) != test.hasError || string(code) != test.expected {
			t.Errorf(`dest("%s") = %q, %v, want %q, error: %v`, test.input, code, err, test.expected, test.hasError)
		}
	}
}

func TestComp(t *testing.T) {
	tests := []struct {
		input    Mnemonic
		expected string
		hasError bool
	}{
		{"0", "0101010", false},
		{"1", "0111111", false},
		{"-1", "0111010", false},
		{"D", "0001100", false},
		{"A", "0110000", false},
		{"M", "1110000", false},
		{"!D", "0001101", false},
		{"!A", "0110001", false},
		{"!M", "1110001", false},
		{"-D", "0001111", false},
		{"-A", "0110011", false},
		{"-M", "1110011", false},
		{"D+1", "0011111", false},
		{"A+1", "0110111", false},
		{"M+1", "1110111", false},
		{"D-1", "0001110", false},
		{"A-1", "0110010", false},
		{"M-1", "1110010", false},
		{"D+A", "0000010", false},
		{"D+M", "1000010", false},
		{"D-A", "0010011", false},
		{"D-M", "1010011", false},
		{"A-D", "0000111", false},
		{"M-D", "1000111", false},
		{"D&A", "0000000", false},
		{"D&M", "1000000", false},
		{"D|A", "0010101", false},
		{"D|M", "1010101", false},
		{"INVALID", "", true}, // Invalid mnemonic
	}

	for _, test := range tests {
		code, err := comp(test.input)
		if (err != nil) != test.hasError || string(code) != test.expected {
			t.Errorf(`comp("%s") = %q, %v, want %q, error: %v`, test.input, code, err, test.expected, test.hasError)
		}
	}
}

func TestJump(t *testing.T) {
	tests := []struct {
		input    Mnemonic
		expected string
		hasError bool
	}{
		{"null", "000", false},
		{"JGT", "001", false},
		{"JEQ", "010", false},
		{"JGE", "011", false},
		{"JLT", "100", false},
		{"JNE", "101", false},
		{"JLE", "110", false},
		{"JMP", "111", false},
		{"INVALID", "", true}, // Invalid mnemonic
	}

	for _, test := range tests {
		code, err := jump(test.input)
		if (err != nil) != test.hasError || string(code) != test.expected {
			t.Errorf(`jump("%s") = %q, %v, want %q, error: %v`, test.input, code, err, test.expected, test.hasError)
		}
	}
}
