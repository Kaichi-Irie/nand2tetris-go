package main

import (
	"nand2tetris-golang/jackcompiler/jackanalyzer"
	"os"
)

func main() {
	err := jackanalyzer.Analize(os.Args[1])
	if err != nil {
		panic(err)
	}
}
