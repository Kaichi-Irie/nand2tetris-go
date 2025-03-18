package vmtranslator

import (
	"errors"
	"fmt"
	"os"
)

func VMTranslator() error {
	fmt.Println("VMTranslator")

	fileName := os.Args[1]
	if fileName[len(fileName)-3:] != ".vm" {
		return errors.New("invalid file extension")
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	parser := NewParser(file, "//")
	codeWriter := NewCodeWriter(fileName)
	defer codeWriter.Close()

	for parser.advance() {
		err := codeWriter.WriteCommand(parser.currentCommand)
		if err != nil {
			return err
		}
	}

	// write infinite loop at the end of the file
	codeWriter.WriteInfinityLoop()

	fmt.Println("done")
	return nil
}
