package symboltable

import "testing"

func TestSymbolTable(t *testing.T) {
	st := NewSymbolTable()

	tests := []struct {
		name          string
		T             string
		kind          string
		expectedName  string
		expectedKind  string
		expectedT     string
		expectedIndex int
	}{
		{"x", "int", VAR, "x", VAR, "int", 0},
		{"y", "int", VAR, "y", VAR, "int", 1},
		{"z", "int", VAR, "z", VAR, "int", 2},
		{"a", "int", FIELD, "a", FIELD, "int", 0},
		{"b", "int", FIELD, "b", FIELD, "int", 1},
		{"c", "int", STATIC, "c", STATIC, "int", 0},
		{"d", "int", STATIC, "d", STATIC, "int", 1},
		{"e", "int", ARG, "e", ARG, "int", 0},
		{"f", "int", ARG, "f", ARG, "int", 1},
	}
	for _, test := range tests {
		err := st.Define(test.name, test.T, test.kind)
		if err != nil {
			t.Errorf("Error defining symbol: %v", err)
		}
		id, exists := st.Lookup(test.name)
		if !exists {
			t.Errorf("Expected to find symbol '%s'", test.name)
		}
		if id.Name != test.expectedName || id.Kind != test.expectedKind || id.T != test.expectedT || id.Index != test.expectedIndex {
			t.Errorf("Expected symbol '%s' to have kind %s, but got %s", test.name, test.expectedKind, id.Kind)
		}
	}
}

func TestReset(t *testing.T) {
	entriesBeforeReset := []struct {
		name string
		T    string
		kind string
	}{
		{"x", "int", VAR},
		{"y", "int", VAR},
		{"z", "int", VAR},
		{"a", "int", FIELD},
		{"b", "int", FIELD},
		{"c", "int", STATIC},
		{"d", "int", STATIC},
		{"e", "int", ARG},
		{"f", "int", ARG},
	}

	entriesAfterReset := []struct {
		name          string
		T             string
		kind          string
		expectedName  string
		expectedKind  string
		expectedT     string
		expectedIndex int
	}{
		{"x", "int", VAR, "x", VAR, "int", 0},
		{"y", "int", VAR, "y", VAR, "int", 1}}

	st := NewSymbolTable()
	for _, entry := range entriesBeforeReset {
		err := st.Define(entry.name, entry.T, entry.kind)
		if err != nil {
			t.Errorf("Error defining symbol: %v", err)
		}
	}
	err := st.Reset()
	if err != nil {
		t.Errorf("Error resetting symbol table: %v", err)
	}

	for _, entry := range entriesAfterReset {
		err := st.Define(entry.name, entry.T, entry.kind)
		if err != nil {
			t.Errorf("Error defining symbol: %v", err)
		}
		id, exists := st.Lookup(entry.name)
		if !exists {
			t.Errorf("Expected to find symbol '%s'", entry.name)
		}
		if id.Name != entry.expectedName || id.Kind != entry.expectedKind || id.T != entry.expectedT || id.Index != entry.expectedIndex {
			t.Errorf("Expected symbol '%s' to have kind %s, but got %s", entry.name, entry.expectedKind, id.Kind)
		}
	}
}
