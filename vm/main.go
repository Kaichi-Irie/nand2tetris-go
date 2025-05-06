package main

import (
	"os"

	"github.com/Kaichi-Irie/nand2tetris-go/vm/vmtranslator"
)

func main() {
	err := vmtranslator.VMTranslator(os.Args[1])
	if err != nil {
		panic(err)
	}
}
