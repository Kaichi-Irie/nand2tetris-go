package main

import (
	"os"

	"github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/jackanalyzer"
)

func main() {
	err := jackanalyzer.Analize(os.Args[1])
	if err != nil {
		panic(err)
	}
}
