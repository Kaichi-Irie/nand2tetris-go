package main

import "nand2tetris-golang/assembler/hack"

func main() {
	err := hack.Hack()
	if err != nil {
		panic(err)
	}
}
