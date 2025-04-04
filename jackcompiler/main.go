package main

import (
	"nand2tetris-golang/jackcompiler/analyzer"
	"os"
)

func main() {
	err := analyzer.Analize(os.Args[1])
	if err != nil {
		panic(err)
	}
}
