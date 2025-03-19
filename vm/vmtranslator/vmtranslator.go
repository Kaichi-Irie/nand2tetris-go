package vmtranslator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func VMTranslator() error {
	fmt.Println("VMTranslator")

	vmFilePath := os.Args[1]
	if filepath.Ext(vmFilePath) != ".vm" {
		return errors.New("invalid file extension")
	}
	file, err := os.Open(vmFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	parser := NewParser(file, "//")
	codeWriter := NewCodeWriter(vmFilePath)
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
