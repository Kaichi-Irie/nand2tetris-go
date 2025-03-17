package hack

/*
Description: `parser.go` は、アセンブリ言語の命令をパースするための構造体と関数を提供する。
- Types:
| types | description |
|-------|-------------|
| Instruction | 命令を表す文字列。例: "@2", "D=M", "(LOOP)" |
| Mnemonic | C命令のdest, comp, jumpのニーモニックを表す文字列。例: "D", "A", "M", "D+1", "A-1", "D-A", "D&A", "JGT" |
| SymbolOrConstant | シンボルまたは定数を表す文字列。例: `"LOOP"`, `"100"` ,`"x"` |
| InstructionType | 命令のタイプを表す列挙型。`A_Instruction`, `C_Instruction`, `L_Instruction` のいずれか。 |

- Functions and methods:
| functions/methods | arguments | return values | description |
|-----------|-----------|---------------|-------------|
| (p *Parser) advance   |           | bool          | `Parser` が次の命令を読み込み、それを現在の命令にする。もしもう命令がない場合は false を返す。 コメント行や空行は無視して，命令が見つかる，またはファイルの終わりに達するまでスキャンを続ける。 |
| (p *Parser) text      |           | string        | 現在の命令を返す。 |
| (p *Parser) symbol    |           | SymbolOrConstant, error | 現在の命令が`"@const"`の場合は`"const"`を返す。`"@variable"`の場合は`"variable"`を返す。`"(label)"`の場合は`"label"`を返す。 A命令でもL命令でもない場合はエラーを返す。 |
| (p *Parser) dest      |           | Mnemonic, error | 現在のC命令のdestニーモニックを返す。 |
| (p *Parser) comp      |           | Mnemonic, error | 現在のC命令のcompニーモニックを返す。 |
| (p *Parser) jump      |           | Mnemonic, error | 現在のC命令のjumpニーモニックを返す。 |
| isConst   | SymbolOrConstant | bool | 引数が定数かどうかを返す。 `string` を `int` に変換できるかどうかで判断する。 |
| getInstructionType | Instruction | InstructionType | 命令のタイプを返す。 |
| isEmptyLine | string | bool | 引数が空行かどうかを返す。 |
| isCommentLine | string | bool | 引数がコメント行かどうかを返す。 |

*/

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

// Instruction is a string that represents an instruction
// Example: "@2", "D=M", "(LOOP)"
type Instruction string

// Mnemonic is a string that represents a mnemonic
// Example: "D", "A", "M", "D+1", "A-1", "D-A", "D&A", "JGT"
type Mnemonic string

// SymbolOrConstant is a string that represents a symbol or a constant
// Example: "LOOP", "100"
type SymbolOrConstant string

// InstructionType is an enum that represents the type of an instruction.
// A_Instruction: @value
// C_Instruction: dest=comp;jump
// L_Instruction: (label)
type InstructionType int

const (
	A_Instruction InstructionType = iota
	C_Instruction
	L_Instruction
)

type Parser struct {
	scanner            *bufio.Scanner
	currentInstruction Instruction
	currentType        InstructionType
}

// advance reads the next instruction from the input and makes it the current instruction
// It returns false if there are no more instructions
// advance ignores empty lines and comments
func (p *Parser) advance() bool {
	ok := p.scanner.Scan()
	if !ok {
		return false
	}
	line := p.text()
	if isEmptyLine(line) || isCommentLine(line) {
		// skip empty or comment line
		return p.advance()
	}
	inst := Instruction(line)
	p.currentInstruction = inst
	p.currentType = getInstructionType(inst)
	return true
}

// text returns the current instruction.
// It removes all spaces from the instruction
func (p *Parser) text() string {
	// remove spaces
	return strings.ReplaceAll(p.scanner.Text(), " ", "")
}

func isEmptyLine(line string) bool {
	return len(line) == 0
}

func isCommentLine(line string) bool {
	return line[0:2] == "//"
}

// isConst returns true if the given SymbolOrConstant is a constant
func isConst(s SymbolOrConstant) bool {
	_, err := strconv.Atoi(string(s))
	return err == nil
}

func getInstructionType(inst Instruction) InstructionType {
	switch {
	case inst[0] == '@':
		return A_Instruction
	case inst[0] == '(' && inst[len(inst)-1] == ')':
		return L_Instruction
	default:
		return C_Instruction
	}
}

// symbol return the symbol or constant of the current instruction
// @100 -> "100"
// @LOOP -> "LOOP"
// (LOOP) -> "LOOP"
func (p *Parser) symbol() (SymbolOrConstant, error) {
	inst := p.currentInstruction
	switch p.currentType {
	case A_Instruction:
		// remove the @
		return SymbolOrConstant(inst[1:]), nil
	case L_Instruction:
		// remove the ( and )
		return SymbolOrConstant(inst[1 : len(inst)-1]), nil
	case C_Instruction:
		return "", errors.New("not an A instruction or L instruction")
	default:
		return "", errors.New("invalid instruction type")
	}

}

// dest return the dest mnemonic of the current C instruction
func (p *Parser) dest() (Mnemonic, error) {
	inst := p.currentInstruction
	switch p.currentType {
	case C_Instruction:
		if i := strings.Index(string(inst), "="); i != -1 {
			return Mnemonic(inst[:i]), nil
		}
		return "null", nil

	case A_Instruction, L_Instruction:
		return "", errors.New("not a C instruction")
	default:
		return "", errors.New("invalid instruction type")
	}
}

func (p *Parser) comp() (Mnemonic, error) {
	inst := p.currentInstruction
	switch p.currentType {
	case C_Instruction:
		// C instruction : dest=comp;jump
		if i := strings.Index(string(inst), "="); i != -1 {
			inst = inst[i+1:]
		}
		if i := strings.Index(string(inst), ";"); i != -1 {
			inst = inst[:i]
		}
		return Mnemonic(inst), nil
	case A_Instruction, L_Instruction:
		return "", errors.New("not a C instruction")
	default:
		return "", errors.New("invalid instruction type")
	}
}

func (p *Parser) jump() (Mnemonic, error) {
	inst := p.currentInstruction
	switch p.currentType {
	case C_Instruction:
		if i := strings.Index(string(inst), ";"); i != -1 {
			return Mnemonic(inst[i+1:]), nil
		}
		return "null", nil
	case A_Instruction, L_Instruction:
		return "", errors.New("not a C instruction")
	default:
		return "", errors.New("invalid instruction type")
	}
}
