package vmtranslator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type CodeWriter struct {
	file         io.WriteCloser
	vmFileStem   string // the base name of the .vm file without the .vm extension. e.g. "SimpleAdd"
	commandCount int    // for generating unique labels
}

// NewCodeWriter creates a new asm file with the given path and returns a CodeWriter. CodeWriter.FileNameStem is set to "", so it must be set before calling WriteCommand.
func NewCodeWriter(asmFilePath string) *CodeWriter {
	if filepath.Ext(asmFilePath) != ".asm" {
		panic("invalid file extension")
	}
	asmFile, err := os.Create(asmFilePath)
	if err != nil {
		panic(err)
	}
	return &CodeWriter{asmFile, "", 0}
}

func (cw *CodeWriter) Write(b []byte) (int, error) {
	return cw.file.Write(b)
}

func (cw *CodeWriter) Close() error {
	return cw.file.Close()
}

// push the value in D register to the stack; RAM[SP]=D, SP++
var push_D = "@SP\nA=M\nM=D\n@SP\nM=M+1\n"

// pop the value from the stack to D register; SP--, D=RAM[SP]
var pop_D = "@SP\nM=M-1\nA=M\nD=M\n"

// pop the value from the stack to RAM[13]; SP--, RAM[13]=RAM[SP]
var pop_R13 = pop_D + "@R13\nM=D\n"

func TranslatePushPop(ctype VMCommandType, seg string, idx int, fileName string) (string, error) {
	var asmcommand string
	THISorTHAT := map[int]string{
		0: "THIS", 1: "THAT",
	}
	// process the segment
	// This process is common to both push and pop commands except for the static segment
	switch seg {
	case "local":
		asmcommand = "@LCL\nD=M\n"
		asmcommand += fmt.Sprintf("@%d\n", idx)
	case "argument":
		asmcommand = "@ARG\nD=M\n"
		asmcommand += fmt.Sprintf("@%d\n", idx)
	case "this":
		asmcommand = "@THIS\nD=M\n"
		asmcommand += fmt.Sprintf("@%d\n", idx)
	case "that":
		asmcommand = "@THAT\nD=M\n"
		asmcommand += fmt.Sprintf("@%d\n", idx)
	case "temp":
		asmcommand = "@5\nD=A\n"
		asmcommand += fmt.Sprintf("@%d\n", idx)
	case "static", "pointer":
		// do nothing
		asmcommand = ""
	case "constant":
		asmcommand = fmt.Sprintf("@%d\n", idx)
	default:
		return "", fmt.Errorf("invalid segment %s", seg)
	}

	// push or pop the value
	switch {
	case ctype == C_PUSH && seg == "constant":
		// push the value to the stack
		// D=idx
		asmcommand += "D=A\n"
		// RAM[SP]=D, SP++ (push D)
		asmcommand += push_D
	case ctype == C_PUSH && seg == "static":
		// push the value to the stack
		// D=RAM[fileName.idx]
		asmcommand += fmt.Sprintf("@%s.%d\nD=M\n", fileName, idx)
		// RAM[SP]=D, SP++ (push D)
		asmcommand += push_D

	case ctype == C_PUSH && seg == "pointer":
		// push the value to the stack
		// D=RAM[THIS(3) or THAT(4)]
		asmcommand += fmt.Sprintf("@%s\nD=M\n", THISorTHAT[idx])
		// RAM[SP]=D, SP++ (push D)
		asmcommand += push_D
	case ctype == C_PUSH:
		// push the value to the stack
		// D=RAM[segbase+idx]
		asmcommand += "A=D+A\nD=M\n"
		// RAM[SP]=D, SP++ (push D)
		asmcommand += push_D
	case ctype == C_POP && seg == "constant":
		return "", fmt.Errorf("cannot pop to constant segment")
	case ctype == C_POP && seg == "static":
		// pop the value from the stack
		// SP--, D=RAM[SP]
		asmcommand += pop_D
		// RAM[fileName.idx]=D
		asmcommand += fmt.Sprintf("@%s.%d\nM=D\n", fileName, idx)
	case ctype == C_POP && seg == "pointer":
		// pop the value from the stack
		// SP--, D=RAM[SP]
		asmcommand += pop_D
		// RAM[3]=D
		asmcommand += fmt.Sprintf("@%s\nM=D\n", THISorTHAT[idx])
	case ctype == C_POP:
		// pop the value from the stack
		// D=segbase+idx, RAM[13]=D
		asmcommand += "D=D+A\n"
		asmcommand += "@R13\nM=D\n"
		// SP--, D=RAM[SP]
		asmcommand += pop_D
		// RAM[segbase+idx]=D
		asmcommand += "@R13\nA=M\nM=D\n"
	default:
		return "", fmt.Errorf("invalid command type %d", ctype)
	}
	return asmcommand, nil
}

func TranslateArithmetic(command VMCommand, cnt int) (string, error) {
	asmcommand := ""
	op := map[VMCommand]string{
		"add": "+", "sub": "-", "and": "&", "or": "|", "neg": "-", "not": "!", "eq": "JEQ", "gt": "JGT", "lt": "JLT"}[command]
	switch command {
	case "add", "sub", "and", "or":
		asmcommand += pop_R13                           // y
		asmcommand += pop_D                             // x
		asmcommand += fmt.Sprintf("@R13\nD=D%sM\n", op) // x op y
		asmcommand += push_D                            // push (x op y)
	case "neg", "not":
		asmcommand += pop_D                      // x
		asmcommand += fmt.Sprintf("D=%sD\n", op) // op x
		asmcommand += push_D                     // push (op x)
	case "eq", "gt", "lt":
		// generate unique labels for the true and false cases
		// Example: EQ_0_TRUE, GT_1_FALSE, LT_2_END, etc.
		prefix := op[1:] + fmt.Sprintf("_%d", cnt)
		asmcommand += pop_R13         // y
		asmcommand += pop_D           // x
		asmcommand += "@R13\nD=D-M\n" // x-y
		// if (x-y op 0) is true, then goto true else goto false
		asmcommand += fmt.Sprintf("@%s_TRUE\nD;%s\n", prefix, op)
		// false case D=0(false)
		asmcommand += fmt.Sprintf("(%s_FALSE)\nD=0\n@%s_END\n0;JMP\n", prefix, prefix)
		// true case D=-1(true)
		asmcommand += fmt.Sprintf("(%s_TRUE)\nD=-1\n@%s_END\n0;JMP\n", prefix, prefix)
		asmcommand += fmt.Sprintf("(%s_END)\n", prefix)
		asmcommand += push_D // push true or false
	default:
		return "", fmt.Errorf("invalid arithmetic command %s", command)
	}
	return asmcommand, nil
}

func TranslateLabel(label string) (string, error) {
	return fmt.Sprintf("(%s)\n", label), nil
}

func TranslateGoto(label string) (string, error) {
	return fmt.Sprintf("@%s\n0;JMP\n", label), nil
}

func TranslateIf(label string) (string, error) {
	asmcommand := pop_D
	asmcommand += fmt.Sprintf("@%s\nD;JNE\n", label)
	return asmcommand, nil
}

// TODO: implemt these.
// func TranslateFunction(label string) (string, error)
// func TranslateCall(label string) (string, error)
// func TranslateReturn(label string) (string, error)

func (cw *CodeWriter) WriteCommand(command VMCommand) error {
	if cw.vmFileStem == "" {
		return fmt.Errorf("fileNameStem is not set")
	}
	// output the command as a comment
	io.WriteString(cw, "// "+string(command)+"\n")
	ctype := getCommandType(command)
	switch ctype {
	case C_ARITHMETIC:
		asmcommand, err := TranslateArithmetic(command, cw.commandCount)
		if err != nil {
			return err
		}
		cw.commandCount++
		_, err = io.WriteString(cw, asmcommand)
		return err
	case C_PUSH, C_POP:
		asmcommand, err := TranslatePushPop(ctype, arg1(command), arg2(command), cw.vmFileStem)
		if err != nil {
			return err
		}
		_, err = io.WriteString(cw, asmcommand)
		return err
	case C_LABEL:
		asmcommand, err := TranslateLabel(arg1(command))
		if err != nil {
			return err
		}
		_, err = io.WriteString(cw, asmcommand)
		return err
	case C_GOTO:
		asmcommand, err := TranslateGoto(arg1(command))
		if err != nil {
			return err
		}
		_, err = io.WriteString(cw, asmcommand)
		return err
	case C_IF:
		asmcommand, err := TranslateIf(arg1(command))
		if err != nil {
			return err
		}
		_, err = io.WriteString(cw, asmcommand)
		return err
	default:
		return fmt.Errorf("invalid command type %d", ctype)
	}
}

func (cw *CodeWriter) WriteInfinityLoop() error {
	// TODO: Avoid label name collision
	label := "INFINITE_LOOP_END"
	_, err := io.WriteString(cw, fmt.Sprintf("// infinite loop\n(%s)\n@%s\n0;JMP\n", label, label))
	return err
}

// TODO: implement bootstrap
// func (cw *CodeWriter) WriteBootStrap() error
