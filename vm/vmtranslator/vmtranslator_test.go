package vmtranslator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

type MyWriteCloser struct {
	io.Writer
}

func (mwc *MyWriteCloser) Close() error {
	// Noop
	return nil
}

// singleTest tests the translation of a single .vm file to .asm file. It reads the .vm file, translates the commands to .asm, and compares the result with the expected .asm file.
func singleTest(vmFilePath string, asmFilePath string) error {
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
	mwc := &MyWriteCloser{buf}
	cw := &CodeWriter{
		file:        mwc,
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
func TestVMTranslation(t *testing.T) {
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
		err := singleTest(test.vmFilePath, test.asmFilePath)
		if err != nil {
			t.Errorf("singleTest failed for %s: %v", test.vmFilePath, err)
		}
	}
}
