package main

import (
	"nand2tetris-golang/vm/vmtranslator"
)

func main() {
	err := vmtranslator.VMTranslator()
	if err != nil {
		panic(err)
	}
}
