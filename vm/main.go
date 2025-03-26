package main

import (
	"nand2tetris-golang/vm/vmtranslator"
	"os"
)

func main() {
	err := vmtranslator.VMTranslator(os.Args[1])
	if err != nil {
		panic(err)
	}
}
