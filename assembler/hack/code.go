package hack

/*
Description: code.go は、アセンブリ言語の命令をバイナリコードに変換するための関数を提供する。

- Types:
| types | description |
|-------|-------------|
| BinaryCode | 15ビットのバイナリコードを表す文字列。例: "000000001100100" |

- Functions:
| functions | arguments | return values | description |
|-----------|-----------|---------------|-------------|
| decimalToBinary | string | BinaryCode, error | `string`で表された10進数を15ビットのバイナリコードに変換する。例: "100" -> "000000001100100" |
| symbol | SymbolOrConstant, *SymbolTable | BinaryCode, error | シンボルまたは定数を15ビットのバイナリコードに変換する。新しい変数が見つかった場合は、シンボルテーブルに追加する。例: "100" -> "000000001100100", "LOOP" -> symbolTable.table["LOOP"]->"000000000000001" |
| dest | Mnemonic | BinaryCode, error | destニーモニックのバイナリコードを返す。 |
| comp | Mnemonic | BinaryCode, error | compニーモニックのバイナリコードを返す。 |
| jump | Mnemonic | BinaryCode, error | jumpニーモニックのバイナリコードを返す。 |

*/

import (
	"errors"
	"strconv"
)

// BinaryCode is a string that represents a 15-bit binary number
type BinaryCode string

// decimalToBinary converts a decimal number to a 15-bit binary number. Example: "100" -> "000000001100100"
func decimalToBinary(decimalExp string) (BinaryCode, error) {
	i, err := strconv.Atoi(decimalExp)
	if err != nil {
		return "", err
	}
	n_digits := 15
	bin_exp := strconv.FormatInt(int64(i), 2)
	for len(bin_exp) < n_digits {
		bin_exp = "0" + bin_exp
	}
	return BinaryCode(bin_exp), nil
}

// symbol converts a symbol or constant to a 15-bit binary number. if new variable is found, add it to the symbol table. Example: "100" -> "000000001100100", "LOOP" -> symbolTable.table["LOOP"]->"000000000000001"
func symbol(symOrConst SymbolOrConstant, symbolTable *SymbolTable) (BinaryCode, error) {
	table, count := symbolTable.table, symbolTable.variableCount
	// if it is a constant
	if isConst(symOrConst) {
		return decimalToBinary(string(symOrConst))
	} else if val, ok := table[symOrConst]; ok {
		// if it is a registered symbol
		return decimalToBinary(strconv.Itoa(val))
	} else {
		// add a new variable to the symbol table
		symbolTable.addVariable(symOrConst)
		return decimalToBinary(strconv.Itoa(count))
	}
}

// dest return the binary code of the dest mnemonic
func dest(mnemonic Mnemonic) (BinaryCode, error) {
	switch mnemonic {
	case "null":
		return "000", nil
	case "M":
		return "001", nil
	case "D":
		return "010", nil
	case "MD":
		return "011", nil
	case "A":
		return "100", nil
	case "AM":
		return "101", nil
	case "AD":
		return "110", nil
	case "AMD":
		return "111", nil
	}
	return "", errors.New("invalid dest mnemonic")
}

// comp return the binary code of the comp mnemonic
func comp(mnemonic Mnemonic) (BinaryCode, error) {
	a := ""
	cccccc := ""
	switch mnemonic {
	case "0", "1", "-1", "D", "A", "!D", "!A", "-D", "-A", "D+1", "A+1", "D-1", "A-1", "D+A", "D-A", "A-D", "D&A", "D|A":
		a = "0"
	case "M", "!M", "-M", "M+1", "M-1", "D+M", "D-M", "M-D", "D&M", "D|M":
		a = "1"
	default:
		return "", errors.New("invalid comp mnemonic")
	}
	switch mnemonic {
	case "0":
		cccccc = "101010"
	case "1":
		cccccc = "111111"
	case "-1":
		cccccc = "111010"
	case "D":
		cccccc = "001100"
	case "A", "M":
		cccccc = "110000"
	case "!D":
		cccccc = "001101"
	case "!A", "!M":
		cccccc = "110001"
	case "-D":
		cccccc = "001111"
	case "-A", "-M":
		cccccc = "110011"
	case "D+1":
		cccccc = "011111"
	case "A+1", "M+1":
		cccccc = "110111"
	case "D-1":
		cccccc = "001110"
	case "A-1", "M-1":
		cccccc = "110010"
	case "D+A", "D+M":
		cccccc = "000010"
	case "D-A", "D-M":
		cccccc = "010011"
	case "A-D", "M-D":
		cccccc = "000111"
	case "D&A", "D&M":
		cccccc = "000000"
	case "D|A", "D|M":
		cccccc = "010101"
	default:
		return "", errors.New("invalid comp mnemonic")
	}
	return BinaryCode(a + cccccc), nil
}

// jump return the binary code of the jump mnemonic
func jump(mnemonic Mnemonic) (BinaryCode, error) {
	switch mnemonic {
	case "null":
		return "000", nil
	case "JGT":
		return "001", nil
	case "JEQ":
		return "010", nil
	case "JGE":
		return "011", nil
	case "JLT":
		return "100", nil
	case "JNE":
		return "101", nil
	case "JLE":
		return "110", nil
	case "JMP":
		return "111", nil
	default:
		return "", errors.New("invalid jump mnemonic")
	}
}
