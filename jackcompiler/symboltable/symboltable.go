package symboltable

import (
	"fmt"
)

const (
	STATIC = "static"
	FIELD  = "field"
	ARG    = "arg"
	VAR    = "var"
	NONE   = "none"
)

var Kinds = []string{
	STATIC,
	FIELD,
	ARG,
	VAR,
	NONE,
}

type Identifier struct {
	Name  string
	Kind  string
	T     string
	Index int
}

type SymbolTable struct {
	StaticCnt   int
	FieldCnt    int
	ArgCnt      int
	VarCnt      int
	VariableMap map[string]Identifier
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		StaticCnt:   0,
		FieldCnt:    0,
		ArgCnt:      0,
		VarCnt:      0,
		VariableMap: make(map[string]Identifier),
	}
}

// Reset resets the symbol table to its initial state.
func (s *SymbolTable) Reset() error {
	s.StaticCnt = 0
	s.FieldCnt = 0
	s.ArgCnt = 0
	s.VarCnt = 0
	s.VariableMap = make(map[string]Identifier)
	return nil
}

/*
Lookup checks if the identifier with the given name exists in the symbol table.
If it exists, it returns the identifier and true. Otherwise, it returns an empty identifier and false.
name: the name of the identifier
*/
func (s *SymbolTable) Lookup(name string) (Identifier, bool) {
	if identifier, exists := s.VariableMap[name]; exists {
		return identifier, true
	}
	return Identifier{}, false
}

/*
Define adds a new identifier to the symbol table with the given name, type, and kind.
name: the name of the identifier
kind: static, field, arg, var
T: int, char, boolean, className
*/
func (s *SymbolTable) Define(name string, T string, kind string) error {
	var index int
	switch kind {
	case STATIC:
		index = s.StaticCnt
		s.StaticCnt++
	case FIELD:
		index = s.FieldCnt
		s.FieldCnt++
	case ARG:
		index = s.ArgCnt
		s.ArgCnt++
	case VAR:
		index = s.VarCnt
		s.VarCnt++
	case NONE:
	default:
		return fmt.Errorf("invalid kind: %s", kind)
	}

	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if T == "" {
		return fmt.Errorf("type cannot be empty")
	}
	if _, exists := s.VariableMap[name]; exists {
		return fmt.Errorf("name %s already defined", name)
	}
	s.VariableMap[name] = Identifier{
		Name:  name,
		Kind:  kind,
		T:     T,
		Index: index,
	}
	fmt.Printf("Added identifier: %s, kind: %s, type: %s, index: %d\n", name, kind, T, index)
	return nil
}
