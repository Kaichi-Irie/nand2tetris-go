package vmwriter

import (
	"fmt"
	"io"
	"slices"
)

const (
	CONSTANT string = "constant"
	ARGUMENT string = "argument"
	LOCAL    string = "local"
	STATIC   string = "static"
	THIS     string = "this"
	THAT     string = "that"
	POINTER  string = "pointer"
	TEMP     string = "temp"
)

var Segments = []string{
	CONSTANT,
	ARGUMENT,
	LOCAL,
	STATIC,
	THIS,
	THAT,
	POINTER,
	TEMP,
}

const (
	ADD string = "add"
	SUB string = "sub"
	NEG string = "neg"
	EQ  string = "eq"
	GT  string = "gt"
	LT  string = "lt"
	AND string = "and"
	OR  string = "or"
	NOT string = "not"
)

var Commands = []string{
	ADD,
	SUB,
	NEG,
	EQ,
	GT,
	LT,
	AND,
	OR,
	NOT,
}

type VMWriter struct {
	w io.WriteCloser
}

func NewVMWriter() *VMWriter {
	return &VMWriter{}
}
func (vmw *VMWriter) writePush(segment string, index int) error {
	if !slices.Contains(Segments, segment) {
		return fmt.Errorf("invalid segment: %s", segment)
	}
	if index < 0 {
		return fmt.Errorf("invalid index: %d", index)
	}
	_, err := fmt.Fprintf(vmw.w, "push %s %d\n", segment, index)
	if err != nil {
		return fmt.Errorf("error writing push command: %w", err)
	}
	return nil
}

func (vmw *VMWriter) writePop(segment string, index int) error {
	if !slices.Contains(Segments, segment) {
		return fmt.Errorf("invalid segment: %s", segment)
	}
	if index < 0 {
		return fmt.Errorf("invalid index: %d", index)
	}
	_, err := fmt.Fprintf(vmw.w, "pop %s %d\n", segment, index)
	if err != nil {
		return fmt.Errorf("error writing pop command: %w", err)
	}
	return nil
}

func (vmw *VMWriter) writeArithmetic(command string) error {
	if !slices.Contains(Commands, command) {
		return fmt.Errorf("invalid command: %s", command)
	}
	_, err := fmt.Fprintf(vmw.w, "%s\n", command)
	if err != nil {
		return fmt.Errorf("error writing arithmetic command: %w", err)
	}
	return nil
}

func (vmw *VMWriter) writeLabel(label string) error {
	_, err := fmt.Fprintf(vmw.w, "label %s\n", label)
	if err != nil {
		return fmt.Errorf("error writing label: %w", err)
	}
	return nil
}

func (vmw *VMWriter) writeGoto(label string) error {
	_, err := fmt.Fprintf(vmw.w, "goto %s\n", label)
	if err != nil {
		return fmt.Errorf("error writing goto: %w", err)
	}
	return nil
}
func (vmw *VMWriter) writeIf(label string) error {
	_, err := fmt.Fprintf(vmw.w, "if-goto %s\n", label)
	if err != nil {
		return fmt.Errorf("error writing if-goto: %w", err)
	}
	return nil
}
func (vmw *VMWriter) writeCall(name string, nArgs int) error {
	if nArgs < 0 {
		return fmt.Errorf("invalid number of arguments: %d", nArgs)
	}
	_, err := fmt.Fprintf(vmw.w, "call %s %d\n", name, nArgs)
	if err != nil {
		return fmt.Errorf("error writing call: %w", err)
	}
	return nil
}

func (vmw *VMWriter) writeFunction(name string, nArgs int) error {
	if nArgs < 0 {
		return fmt.Errorf("invalid number of arguments: %d", nArgs)
	}
	_, err := fmt.Fprintf(vmw.w, "function %s %d\n", name, nArgs)
	if err != nil {
		return fmt.Errorf("error writing function: %w", err)
	}
	return nil
}

func (vmw *VMWriter) writeReturn() error {
	_, err := fmt.Fprintf(vmw.w, "return\n")
	if err != nil {
		return fmt.Errorf("error writing return: %w", err)
	}
	return nil
}
func (vmw *VMWriter) close() error {
	if err := vmw.w.Close(); err != nil {
		return fmt.Errorf("error closing writer: %w", err)
	}
	return nil
}
