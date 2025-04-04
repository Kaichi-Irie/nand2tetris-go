package tokenizer

import "testing"

func TestExtractIntConst(t *testing.T) {
	tests := []struct {
		token string
		want  int
	}{
		{"123", 123},
		{"0", 0},
		{"999", 999},
		{"098", 98},
		{"098hello", 98},
		{"3+5", 3},
		{"3   ", 3},
		{"4*2", 4},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got, err := extractIntConst(tt.token); err != nil {
				t.Errorf("extractIntConst(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("extractIntConst(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}
func TestExtractStringConst(t *testing.T) {
	tests := []struct {
		token string
		want  string
	}{
		{"\"hello\"", "hello"},
		{"\"\"", ""},
		{"\"123\"", "123"},
		{"\"hello world\"", "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got, err := extractStringConst(tt.token); err != nil {
				t.Errorf("extractStringConst(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("extractStringConst(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestExtractIdentifier(t *testing.T) {
	tests := []struct {
		token string
		want  string
	}{
		{"hello + 3", "hello"},
		{"_hello(\"str\")", "_hello"},
		{"hello_world", "hello_world"},
		{"hello123", "hello123"},
		{"x)", "x"},
		{"x.y", "x"},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got, err := extractIdentifier(tt.token); err != nil {
				t.Errorf("extractIdentifier(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("extractIdentifier(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}
