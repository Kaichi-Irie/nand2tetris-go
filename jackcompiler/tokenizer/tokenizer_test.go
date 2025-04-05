package tokenizer

import "testing"

func TestIsIntConst(t *testing.T) {
	tests := []struct {
		token string
		want  bool
	}{
		{"123", true},
		{"0", true},
		{"999", true},
		{"098", true},
		{"098hello", false},
		{"3+5", false},
		{"3   ", false},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got := IsIntConst(tt.token); got != tt.want {
				t.Errorf("isIntConst(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}
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
			if got, err := ExtractIntConst(tt.token); err != nil {
				t.Errorf("extractIntConst(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("extractIntConst(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestIsStringConst(t *testing.T) {
	tests := []struct {
		token string
		want  bool
	}{
		{"\"hello\"", true},
		{"\"\"", true},
		{"\"123\"", true},
		{"\"hello world\"", true},
		{"hello", false},
		{"123", false},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got := IsStringConst(tt.token); got != tt.want {
				t.Errorf("isStringConst(%s) = %v, want %v", tt.token, got, tt.want)
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
			if got, err := ExtractStringConst(tt.token); err != nil {
				t.Errorf("extractStringConst(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("extractStringConst(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}

func TestIsIdentifier(t *testing.T) {
	tests := []struct {
		token string
		want  bool
	}{
		{"hello", true},
		{"_hello", true},
		{"hello_world", true},
		{"hello123", true},
		{"3hello", false},
		{" hello", false},
		{"hello ", false},
		{"hello+world", false},
		{"hello-world", false},
		{"hello.world", false},
		{"hello,world", false},
		{"hello;world", false},
		{"hello:world", false},
		{"hello?world", false},
		{"hello!world", false},
		{"hello@world", false},
		{"hello#world", false},
		{"hello$world", false},
		{"hello%world", false},
		{"hello^world", false},
		{"hello&world", false},
		{"hello*world", false},
		{"hello(world", false},
		{"hello)world", false},
		{"class", false},
		{"method", false},
		{"function", false},
		{"constructor", false},
		{"int", false},
		{"boolean", false},
		{"char", false},
		{"void", false},
		{"var", false},
		{"static", false},
		{"field", false},
		{"let", false},
		{"do", false},
		{"if", false},
		{"else", false},
		{"while", false},
		{"return", false},
		{"true", false},
		{"false", false},
		{"null", false},
		{"this", false},
		{"123", false},
		{"\"hello\"", false},
		{"\"hello world\"", false},
		{"(", false}}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			if got := IsIdentifier(tt.token); got != tt.want {
				t.Errorf("isIdentifier(%s) = %v, want %v", tt.token, got, tt.want)
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
			if got, err := ExtractIdentifier(tt.token); err != nil {
				t.Errorf("extractIdentifier(%s) = %v, want %v", tt.token, err, tt.want)
			} else if got != tt.want {
				t.Errorf("extractIdentifier(%s) = %v, want %v", tt.token, got, tt.want)
			}
		})
	}
}
