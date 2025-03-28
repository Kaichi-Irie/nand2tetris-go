package vmtranslator

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// VMTranslator translates VM code to Hack assembly code. The input can be a .vm file or a directory containing .vm files. The output is a .asm file with the same name as the input file or directory.
func VMTranslator(path string) error {
	// path can be a .vm file or a directory containing .vm files
	fmt.Println("VMTranslator")
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	var asmFilePath string
	var vmFilePaths []string

	if info.IsDir() {
		dirName := filepath.Base(path)
		asmFilePath = filepath.Join(path, dirName+".asm")
		vmFilePaths, err = filepath.Glob(filepath.Join(path, "*.vm"))
		if err != nil {
			panic(err)
		}
		// Sys.vm must be included in the list of .vm files
		if exists := slices.Contains(vmFilePaths, filepath.Join(path, "Sys.vm")); !exists {
			return fmt.Errorf("Sys.vm must be included in the given directory")
		}
	} else if filepath.Ext(path) == ".vm" {
		asmFilePath = path[:len(path)-3] + ".asm"
		vmFilePaths = []string{path}
	} else {
		return fmt.Errorf("input file must be a .vm file or a directory")
	}

	codeWriter := NewCodeWriter(asmFilePath)
	defer codeWriter.Close()
	if info.IsDir() {
		// If the input is a directory, write the bootstrap code at the beginning of the .asm file. The bootstrap code initializes the stack pointer to 256 and calls Sys.init.
		err := codeWriter.WriteBootStrap()
		if err != nil {
			panic(err)
		}
	}

	for _, vmFilePath := range vmFilePaths {
		vmFile, err := os.Open(vmFilePath)
		if err != nil {
			return err
		}
		defer vmFile.Close()

		// vmFileBase is the base name of the .vm file with the .vm extension. e.g. "SimpleAdd.vm"
		vmFileBase := filepath.Base(vmFilePath)
		codeWriter.vmFileStem = vmFileBase[:len(vmFileBase)-3]
		parser := NewParser(vmFile, "//")
		for parser.advance() {
			err := codeWriter.WriteCommand(parser.currentCommand)
			if err != nil {
				return err
			}
		}
	}

	// If the input is a file, write an infinite loop at the end of the .asm file
	if !info.IsDir() {
		codeWriter.WriteInfinityLoop()
	}

	fmt.Println("done")
	return nil
}
