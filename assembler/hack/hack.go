// hack package provides functions to convert Hack assembly language programs to binary code.
package hack

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// Hack converts an assembly language file to a binary code file. The file name is passed as a command line argument. The output file has the same name as the input file but with a .hack extension.
func Hack(fileName string) error {
	fmt.Println("Hack")
	if fileName[len(fileName)-4:] != ".asm" {
		return errors.New("invalid file extension")
	}

	// first pass looks for only L instructions and add them to the symbol table
	symbolTable, err := firstPass(fileName)
	if err != nil {
		return err
	}

	// create a new file with the same name as the input file but with .hack extension
	hackFile, err := os.Create(fileName[:len(fileName)-4] + ".hack")
	if err != nil {
		return err
	}
	defer hackFile.Close()

	// second pass looks for A and C instructions and convert them to binary code
	err = secondPass(fileName, hackFile, symbolTable)
	if err != nil {
		return err
	}
	fmt.Println("done")
	return nil
}

// firstPass looks for L instructions and add them to the symbol table.
func firstPass(fileName string) (SymbolTable, error) {
	fmt.Println("first pass")
	asmFile, err := os.Open(fileName)
	if err != nil {
		return SymbolTable{}, err
	}
	defer asmFile.Close()

	p := Parser{scanner: bufio.NewScanner(asmFile)}
	symbolTable := NewSymbolTable()

	// first pass: build symbol table and add labels to it
	for count := 0; p.advance(); {
		if p.currentType == L_Instruction {
			labelSymbol, err := p.symbol()
			if err != nil {
				return SymbolTable{}, err
			}
			// add label to the symbol table
			symbolTable.addLabel(labelSymbol, count)
		} else if p.currentType == A_Instruction || p.currentType == C_Instruction {
			count++
		} else {
			return SymbolTable{}, errors.New("invalid instruction type")
		}
	}
	return symbolTable, nil
}

// secondPass looks for A and C instructions and convert them to binary code.
func secondPass(fileName string, hackFile *os.File, symbolTable SymbolTable) error {
	fmt.Println("second pass")

	asmFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer asmFile.Close()

	p := Parser{scanner: bufio.NewScanner(asmFile)}
	// second pass
	for p.advance() {
		instType := p.currentType
		switch instType {
		case A_Instruction:
			symOrConst, err := p.symbol()
			if err != nil {
				return err
			}
			symbolCode, err := symbol(symOrConst, &symbolTable)
			if err != nil {
				return err
			}
			code := "0" + string(symbolCode)
			_, err = hackFile.WriteString(code + "\n")
			if err != nil {
				return err
			}
		case C_Instruction:
			destMnemonic, err := p.dest()
			if err != nil {
				return err
			}
			compMnemonic, err := p.comp()
			if err != nil {
				return err
			}
			jumpMnemonic, err := p.jump()
			if err != nil {
				return err
			}
			destCode, err := dest(destMnemonic)
			if err != nil {
				return err
			}
			compCode, err := comp(compMnemonic)
			if err != nil {
				return err
			}
			jumpCode, err := jump(jumpMnemonic)
			if err != nil {
				return err
			}
			code := "111" + string(compCode) + string(destCode) + string(jumpCode)
			_, err = hackFile.WriteString(code + "\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
