package hack

/*
Description:
`hack` パッケージは、Hackコンピュータのアセンブリ言語プログラムをバイナリコードに変換するための関数を提供する。

- Functions:
| functions | arguments | return values | description |
|-----------|-----------|---------------|-------------|
| Hack | string | error | アセンブリ言語のファイルをバイナリコードに変換する。 ファイル名はコマンドライン引数として渡される。 書き込みファイルは入力ファイルと同じ名前で、拡張子が`.hack`になる。 |
| firstPass | string | SymbolTable, error | 1回目のパス。L命令を探し、それらをシンボルテーブルに追加する。 |
| secondPass | string, *os.File, SymbolTable | error | 2回目のパス。A命令とC命令を探し、それらをバイナリコードに変換する。 |

*/

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func Hack() error {
	fmt.Println("Hack")

	fileName := os.Args[1]
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
			labelSymbol, _ := p.symbol()
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
			symOrConst, _ := p.symbol()
			symbolCode, _ := symbol(symOrConst, &symbolTable)
			code := "0" + string(symbolCode)
			_, err := hackFile.WriteString(code + "\n")
			if err != nil {
				return err
			}
		case C_Instruction:
			destMnemonic, _ := p.dest()
			compMnemonic, _ := p.comp()
			jumpMnemonic, _ := p.jump()
			destCode, _ := dest(destMnemonic)
			compCode, _ := comp(compMnemonic)
			jumpCode, _ := jump(jumpMnemonic)
			code := "111" + string(compCode) + string(destCode) + string(jumpCode)
			_, err := hackFile.WriteString(code + "\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
