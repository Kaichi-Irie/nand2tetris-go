package main

import (
	"os"

	"github.com/Kaichi-Irie/nand2tetris-go/assembler/hack"
)

func main() {
	err := hack.Hack(os.Args[1])
	if err != nil {
		panic(err)
	}
}
