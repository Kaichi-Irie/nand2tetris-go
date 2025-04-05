package main

import (
	"nand2tetris-go/vm/vmtranslator"
	"os"
)

func main() {
	err := vmtranslator.VMTranslator(os.Args[1])
	if err != nil {
		panic(err)
	}
}
