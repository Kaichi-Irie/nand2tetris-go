package vmtranslator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// testSingleFile tests the translation of a single .vm file to .asm file. It reads the .vm file, translates the commands to .asm, and compares the result with the expected .asm file.
func testSingleFile(vmFilePath string, asmFilePath string) error {
	if filepath.Ext(vmFilePath) != ".vm" {
		return fmt.Errorf("invalid file extension")
	}

	vmFile, err := os.Open(vmFilePath)
	if err != nil {
		return err
	}
	defer vmFile.Close()

	if filepath.Ext(asmFilePath) != ".asm" {
		return fmt.Errorf("invalid file extension")
	}
	asmFile, err := os.Open(asmFilePath)
	if err != nil {
		return err
	}
	defer asmFile.Close()

	buf := &bytes.Buffer{}
	cw := &CodeWriter{
		file:        buf,
		vmFileStem:  filepath.Base(vmFilePath)[:len(filepath.Base(vmFilePath))-3],
		returnCount: make(map[string]int),
	}

	parser := NewParser(vmFile, "//")
	for parser.advance() {
		err := cw.WriteCommand(parser.currentCommand)
		if err != nil {
			return err
		}
	}
	cw.WriteInfinityLoop()

	buf2 := &bytes.Buffer{}
	io.Copy(buf2, asmFile)

	if expected := buf2.String(); buf.String() != expected {
		return fmt.Errorf("Expected %s, got %s", expected, buf.String())
	}
	return nil
}

// testSingleDir tests the translation of a directory containing .vm files to a single .asm file. It reads the .vm files, translates the commands to .asm, and compares the result with the expected .asm file. The directory must contain a Sys.vm file that is the entry point of the program.
func testSingleDir(dirName string, asmFilePath string) error {
	if info, err := os.Stat(dirName); err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("input must be a directory")
	}

	if filepath.Ext(asmFilePath) != ".asm" {
		return fmt.Errorf("invalid file extension")
	}
	asmFile, err := os.Open(asmFilePath)
	if err != nil {
		return err
	}
	defer asmFile.Close()

	buf := &bytes.Buffer{}
	cw := &CodeWriter{
		file:        buf,
		returnCount: make(map[string]int),
	}

	vmFilePaths, err := filepath.Glob(filepath.Join(dirName, "*.vm"))
	if err != nil {
		return err
	}
	err = cw.WriteBootStrap()
	if err != nil {
		return err
	}
	for _, vmFilePath := range vmFilePaths {
		vmFile, err := os.Open(vmFilePath)
		if err != nil {
			return err
		}
		defer vmFile.Close()

		// vmFileBase is the base name of the .vm file with the .vm extension. e.g. "SimpleAdd.vm"
		vmFileBase := filepath.Base(vmFilePath)
		cw.vmFileStem = vmFileBase[:len(vmFileBase)-3]
		parser := NewParser(vmFile, "//")
		for parser.advance() {
			err := cw.WriteCommand(parser.currentCommand)
			if err != nil {
				return err
			}
		}
	}

	buf2 := &bytes.Buffer{}
	io.Copy(buf2, asmFile)

	if expected := buf2.String(); buf.String() != expected {
		return fmt.Errorf("Expected %s, got %s", expected, buf.String())
	}
	return nil

}
func TestVMTranslationSingleFiles(t *testing.T) {
	tests := []struct {
		vmFilePath  string
		asmFilePath string
	}{
		{"../vm_files/SimpleAdd.vm", "../vm_files/SimpleAdd.asm"},
		{"../vm_files/BasicTest.vm", "../vm_files/BasicTest.asm"},
		{"../vm_files/BasicLoop.vm", "../vm_files/BasicLoop.asm"},
		{"../vm_files/StackTest.vm", "../vm_files/StackTest.asm"},
		{"../vm_files/StaticTest.vm", "../vm_files/StaticTest.asm"},
		{"../vm_files/FibonacciSeries.vm", "../vm_files/FibonacciSeries.asm"},
		{"../vm_files/PointerTest.vm", "../vm_files/PointerTest.asm"},
		{"../vm_files/SimpleFunction.vm", "../vm_files/SimpleFunction.asm"},
		{"../vm_files/NestedCall/Sys.vm", "../vm_files/NestedCall/Sys.asm"},
	}
	for _, test := range tests {
		err := testSingleFile(test.vmFilePath, test.asmFilePath)
		if err != nil {
			t.Errorf("singleTest failed for %s: %v", test.vmFilePath, err)
		}
	}
}

// TestVMTranslationDirs tests the translation of directories containing .vm files to .asm files. It reads the .vm files, translates the commands to .asm, and compares the result with the expected .asm files. The directories must contain a Sys.vm file that is the entry point of the program.
func TestVMTranslationDirs(t *testing.T) {
	tests := []struct {
		dirName     string
		asmFilePath string
	}{
		{"../vm_files/FibonacciElement", "../vm_files/FibonacciElement/FibonacciElement.asm"},
		{"../vm_files/StaticsTest", "../vm_files/StaticsTest/StaticsTest.asm"},
		{"../vm_files/NestedCall", "../vm_files/NestedCall/NestedCall.asm"},
	}
	for _, test := range tests {
		err := testSingleDir(test.dirName, test.asmFilePath)
		if err != nil {
			t.Errorf("dirTest failed for %s: %v", test.dirName, err)
		}
	}
}
