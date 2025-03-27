package main

import (
	"nand2tetris-golang/assembler/hack"
	"os"
)

func main() {
	err := hack.Hack(os.Args[1])
	if err != nil {
		panic(err)
	}
}
