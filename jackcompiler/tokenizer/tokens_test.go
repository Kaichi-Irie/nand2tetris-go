package tokenizer

import "testing"

func TestIsOp(t *testing.T) {
	tests := []struct {
		token  Token
		expect bool
	}{
		{PLUS, true},
		{MINUS, true},
		{ASTERISK, true},
		{SLASH, true},
		{AND, true},
		{OR, true},
		{LESS, true},
		{GREATER, true},
		{EQUAL, true},
		{CLASS, false},
		{METHOD, false},
		{FUNCTION, false},
		{CONSTRUCTOR, false},
		{INT, false},
		{BOOLEAN, false},
		{CHAR, false},
		{VOID, false},
		{VAR, false},
		{STATIC, false},
		{FIELD, false},
		{LET, false},
		{DO, false},
		{IF, false},
		{ELSE, false},
		{WHILE, false},
		{RETURN, false},
		{TRUE, false},
		{FALSE, false},
		{NULL, false},
		{THIS, false},
	}

	for _, tt := range tests {
		t.Run(tt.token.Val, func(t *testing.T) {
			if got := tt.token.IsOp(); got != tt.expect {
				t.Errorf("IsOp() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestIsUnaryOp(t *testing.T) {
	tests := []struct {
		token  Token
		expect bool
	}{
		{NOT, true},
		{MINUS, true},
		{PLUS, false},
		{ASTERISK, false},
		{SLASH, false},
		{AND, false},
		{OR, false},
		{LESS, false},
		{GREATER, false},
		{EQUAL, false},
	}

	for _, tt := range tests {
		t.Run(tt.token.Val, func(t *testing.T) {
			if got := tt.token.IsUnaryOp(); got != tt.expect {
				t.Errorf("IsUnaryOp() = %v, want %v", got, tt.expect)
			}
		})
	}
}
func TestIsKeywordConst(t *testing.T) {
	tests := []struct {
		token  Token
		expect bool
	}{
		{TRUE, true},
		{FALSE, true},
		{NULL, true},
		{THIS, true},
		{INT, false},
		{BOOLEAN, false},
		{CHAR, false},
		{VOID, false},
	}

	for _, tt := range tests {
		t.Run(tt.token.Val, func(t *testing.T) {
			if got := tt.token.IsKeywordConst(); got != tt.expect {
				t.Errorf("IsKeywordConst() = %v, want %v", got, tt.expect)
			}
		})
	}
}
func TestIsPrimitiveType(t *testing.T) {
	tests := []struct {
		token  Token
		expect bool
	}{
		{INT, true},
		{BOOLEAN, true},
		{CHAR, true},
		{VOID, true},
		{CLASS, false},
		{METHOD, false},
		{FUNCTION, false},
		{CONSTRUCTOR, false},
	}

	for _, tt := range tests {
		t.Run(tt.token.Val, func(t *testing.T) {
			if got := tt.token.IsPrimitiveType(); got != tt.expect {
				t.Errorf("IsPrimitiveType() = %v, want %v", got, tt.expect)
			}
		})
	}
}
