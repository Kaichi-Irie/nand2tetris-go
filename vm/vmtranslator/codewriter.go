package vmtranslator

import (
	"fmt"
	"io"
	"os"
)

type CodeWriter struct {
	file         io.WriteCloser
	fileNameStem string
}

func NewCodeWriter(fileName string) *CodeWriter {
	if extension := fileName[len(fileName)-3:]; extension != ".vm" {
		panic("invalid file extension")
	}
	fileNameStem := fileName[:len(fileName)-3]
	file, err := os.Create(fileNameStem + ".asm")
	if err != nil {
		panic(err)
	}
	return &CodeWriter{file, fileNameStem}
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
		asmcommand = "@R5\nD=M\n"
		asmcommand += fmt.Sprintf("@%d\n", idx)
	case "pointer":
		switch idx {
		case 0:
			asmcommand = "@THIS\nD=M\n"
			asmcommand += fmt.Sprintf("@%d\n", idx)
		case 1:
			asmcommand = "@THAT\nD=M\n"
			asmcommand += fmt.Sprintf("@%d\n", idx)
		default:
			return "", fmt.Errorf("invalid pointer index %d", idx)
		}
	case "static":
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

func TranslateArithmetic(command VMCommand) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (cw *CodeWriter) WritePushPop(ctype VMCommandType, seg string, idx int) error {
	if ctype != C_PUSH && ctype != C_POP {
		return fmt.Errorf("invalid command type %d", ctype)
	}

	// translate the command
	asmcommand, err := TranslatePushPop(ctype, seg, idx, cw.fileNameStem)
	if err != nil {
		return err
	}
	_, err = io.WriteString(cw, asmcommand)
	return err
}
func (cw *CodeWriter) WriteArithmetic(command VMCommand) error {
	if ctype := getCommandType(command); ctype != C_ARITHMETIC {
		return fmt.Errorf("invalid command. expected C_ARITHMETIC, got %d", ctype)
	}
	asmcommand, err := TranslateArithmetic(command)
	if err != nil {
		return err
	}
	io.WriteString(cw, asmcommand)
	return nil
}

func (cw *CodeWriter) WriteCommand(command VMCommand) error {
	// output the command as a comment
	io.WriteString(cw, "// "+string(command)+"\n")
	ctype := getCommandType(command)
	switch ctype {
	case C_ARITHMETIC:
		return cw.WriteArithmetic(command)
	case C_PUSH, C_POP:
		return cw.WritePushPop(ctype, arg1(command), arg2(command))
	default:
		return fmt.Errorf("invalid command type %d", ctype)
	}
}

func (cw *CodeWriter) WriteInfinityLoop() error {
	_, err := io.WriteString(cw, "// infinite loop\n(END)\n@END\n0;JMP\n")
	return err
}
