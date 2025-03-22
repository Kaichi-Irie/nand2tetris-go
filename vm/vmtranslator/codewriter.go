package vmtranslator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type CodeWriter struct {
	file         io.WriteCloser
	fileNameStem string
	commandCount int // for generating unique labels
}

func NewCodeWriter(vmFilePath string) *CodeWriter {
	if filepath.Ext(vmFilePath) != ".vm" {
		panic("invalid file extension")
	}
	file, err := os.Create(vmFilePath[:len(vmFilePath)-2] + "asm")
	if err != nil {
		panic(err)
	}
	fileNameStem := vmFilePath[:len(vmFilePath)-3]
	return &CodeWriter{file, fileNameStem, 0}
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
		asmcommand = "@5\nD=A\n"
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

// TODO: implemt these.
func TranslateLabel(label string) (string, error)
func TranslateGoto(label string) (string, error)
func TranslateIf(label string) (string, error)
func TranslateFunction(label string) (string, error)
func TranslateCall(label string) (string, error)
func TranslateReturn(label string) (string, error)

func (cw *CodeWriter) WriteCommand(command VMCommand) error {
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
		asmcommand, err := TranslatePushPop(ctype, arg1(command), arg2(command), cw.fileNameStem)
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
	_, err := io.WriteString(cw, "// infinite loop\n(END)\n@END\n0;JMP\n")
	return err
}

// TODO: implement bootstrap
// func (cw *CodeWriter) WriteBootStrap() error
