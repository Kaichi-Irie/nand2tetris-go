package hack

/*
Description: symbol_table.go は、シンボルテーブルを表す構造体と関数を提供する。

- Types:
| types | description |
|-------|-------------|
| SymbolTable | シンボルテーブルを表す構造体。 変数用のRAM空間を管理するカウンタ`variableCount`と、シンボルテーブルを表す`map[SymbolOrConstant]int`を持つ。 |

- Functions:
| functions | arguments | return values | description |
|-----------|-----------|---------------|-------------|
| NewSymbolTable | | SymbolTable | 新しいシンボルテーブルを返す。 |
| (s *SymbolTable) contains | SymbolOrConstant | bool | シンボルテーブルに引数のシンボルが含まれているかどうかを返す。 |
| (s *SymbolTable) addVariable | SymbolOrConstant | | シンボルテーブルに変数を追加する。 |
| (s *SymbolTable) addLabel | SymbolOrConstant, int | | シンボルテーブルにラベルを追加する。 |


*/

import "errors"

type SymbolTable struct {
	table         map[SymbolOrConstant]int
	variableCount int
}

func NewSymbolTable() SymbolTable {
	table := map[SymbolOrConstant]int{
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15,
		"SCREEN": 16384,
		"KBD":    24576,
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
	}
	return SymbolTable{
		table:         table,
		variableCount: 16,
	}
}

func (s *SymbolTable) contains(symbol SymbolOrConstant) bool {
	_, ok := s.table[symbol]
	return ok
}

// addVariable adds a variable to the symbol table
func (s *SymbolTable) addVariable(symbol SymbolOrConstant) {
	s.table[symbol] = s.variableCount
	s.variableCount++
}

// addLabel adds a label to the symbol table
func (s *SymbolTable) addLabel(label SymbolOrConstant, count int) {
	if !s.contains(label) {
		s.table[label] = count
	} else {
		panic(errors.New("defined label is defined again"))
	}
}
