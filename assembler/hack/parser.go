package hack

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

// Instruction is a string that represents an instruction. Example: "@2", "D=M", "(LOOP)"
type Instruction string

// Mnemonic is a string that represents a mnemonic. Example: "D", "A", "M", "D+1", "A-1", "D-A", "D&A", "JGT"
type Mnemonic string

// SymbolOrConstant is a string that represents a symbol or a constant. Example: "LOOP", "100"
type SymbolOrConstant string

// InstructionType is an enum that represents the type of an instruction. A_Instruction: @value, C_Instruction: dest=comp;jump, L_Instruction: (label)
type InstructionType int

const (
	A_Instruction InstructionType = iota
	C_Instruction
	L_Instruction
)

// Parser is a struct that represents a parser for the Hack assembly language
type Parser struct {
	scanner            *bufio.Scanner
	currentInstruction Instruction
	currentType        InstructionType
}

func NewParser(scanner *bufio.Scanner) *Parser {
	return &Parser{
		scanner: scanner,
	}
}

// advance reads the next instruction from the input and makes it the current instruction. It returns false if there are no more instructions. The function ignores empty lines and comments.
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

// text returns the current instruction. It removes all spaces from the instruction
func (p *Parser) text() string {
	// remove spaces
	return strings.ReplaceAll(p.scanner.Text(), " ", "")
}

// isEmptyLine returns true if the given line is empty
func isEmptyLine(line string) bool {
	return len(line) == 0
}

// isCommentLine returns true if the given line is a comment
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

/*
symbol return the symbol or constant of the current instruction
Example: @100 -> "100", @LOOP -> "LOOP", (LOOP) -> "LOOP"
*/
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

// comp return the comp mnemonic of the current C instruction
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

// jump return the jump mnemonic of the current C instruction
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
