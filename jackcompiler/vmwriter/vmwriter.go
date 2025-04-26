package vmwriter

import (
	"fmt"
	"io"
	st "nand2tetris-go/jackcompiler/symboltable"
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

var SegmentOfKind = map[string]string{
	st.STATIC: STATIC,
	st.FIELD:  THIS,
	st.ARG:    ARGUMENT,
	st.VAR:    LOCAL,
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
	MUL string = "call Math.multiply 2"
	DIV string = "call Math.divide 2"
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
	MUL,
	DIV,
}

type VMWriter struct {
	w io.Writer
}

func New(w io.Writer) *VMWriter {
	return &VMWriter{w: w}
}
func (vmw *VMWriter) WritePush(segment string, index int) error {
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

func (vmw *VMWriter) WritePop(segment string, index int) error {
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

func (vmw *VMWriter) WriteArithmetic(command string) error {
	if !slices.Contains(Commands, command) {
		return fmt.Errorf("invalid command: %s", command)
	}
	_, err := fmt.Fprintf(vmw.w, "%s\n", command)
	if err != nil {
		return fmt.Errorf("error writing arithmetic command: %w", err)
	}
	return nil
}

func (vmw *VMWriter) WriteLabel(label string) error {
	_, err := fmt.Fprintf(vmw.w, "label %s\n", label)
	if err != nil {
		return fmt.Errorf("error writing label: %w", err)
	}
	return nil
}

func (vmw *VMWriter) WriteGoto(label string) error {
	_, err := fmt.Fprintf(vmw.w, "goto %s\n", label)
	if err != nil {
		return fmt.Errorf("error writing goto: %w", err)
	}
	return nil
}
func (vmw *VMWriter) WriteIf(label string) error {
	_, err := fmt.Fprintf(vmw.w, "if-goto %s\n", label)
	if err != nil {
		return fmt.Errorf("error writing if-goto: %w", err)
	}
	return nil
}
func (vmw *VMWriter) WriteCall(name string, nArgs int) error {
	if nArgs < 0 {
		return fmt.Errorf("invalid number of arguments: %d", nArgs)
	}
	_, err := fmt.Fprintf(vmw.w, "call %s %d\n", name, nArgs)
	if err != nil {
		return fmt.Errorf("error writing call: %w", err)
	}
	return nil
}

func (vmw *VMWriter) WriteFunction(name string, nVars int) error {
	if nVars < 0 {
		return fmt.Errorf("invalid number of arguments: %d", nVars)
	}
	_, err := fmt.Fprintf(vmw.w, "function %s %d\n", name, nVars)
	if err != nil {
		return fmt.Errorf("error writing function: %w", err)
	}
	return nil
}

func (vmw *VMWriter) WriteReturn(isVoid bool) error {
	if isVoid {
		err := vmw.WritePush(CONSTANT, 0)
		if err != nil {
			return err
		}
	}
	_, err := io.WriteString(vmw.w, "return\n")
	if err != nil {
		return fmt.Errorf("error writing return: %w", err)
	}
	return nil
}

// func (vmw *VMWriter) Close() error {
// 	if err := vmw.w.Close(); err != nil {
// 		return fmt.Errorf("error closing writer: %w", err)
// 	}
// 	return nil
// }
