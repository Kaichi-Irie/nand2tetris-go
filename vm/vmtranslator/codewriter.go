package vmtranslator

import (
	"fmt"
	"io"
)

type CodeWriter struct {
	file     io.WriteCloser
	fileName string
}

func (cw *CodeWriter) Write(b []byte) (int, error) {
	return cw.file.Write(b)
}

func (cw *CodeWriter) Close() error {
	return cw.file.Close()
}

func TranslatePushPop(ctype VMCommandType, seg string, idx int, fileName string) (string, error) {
	var asmcommand string
	switch seg {
	case "local":
		asmcommand = "@LCL\n"
	case "argument":
		asmcommand = "@ARG\n"
	case "this":
		asmcommand = "@THIS\n"
	case "that":
		asmcommand = "@THAT\n"
	case "temp":
		asmcommand = "@R5\n"
	case "pointer":
		switch idx {
		case 0:
			asmcommand = "@THIS\n"
		case 1:
			asmcommand = "@THAT\n"
		default:
			return "", fmt.Errorf("invalid pointer index %d", idx)
		}
	case "static":
		asmcommand = fmt.Sprintf("@%s.%d\n", fileName, idx)
	case "constant":
		// do nothing
		asmcommand = ""
	default:
		return "", fmt.Errorf("invalid segment %s", seg)
	}
	asmcommand += "D=M\n"
	asmcommand += fmt.Sprintf("@%d\n", idx)

	// push or pop the value
	switch ctype {
	case C_PUSH:
		// push the value to the stack
		// D=idx if segment is constant else D=RAM[segbase+idx]
		if seg == "constant" {
			asmcommand += "D=A\n"
		} else {
			asmcommand += "A=D+A\nD=M\n"
		}
		// RAM[SP] = D, SP++
		asmcommand += "@SP\nA=M\nM=D\n@SP\nM=M+1\n"
	case C_POP:
		// pop the value from the stack
		// D=segbase+idx, RAM[13]=D
		asmcommand += "D=D+A\n"
		asmcommand += "@R13\nM=D\n"
		// SP--, RAM[segbase+idx] = RAM[SP]
		asmcommand += "@SP\nM=M-1\nA=M\nD=M\n@R13\nA=M\nM=D\n"
	default:
		return "", fmt.Errorf("invalid command type %d", ctype)
	}
	return "", nil
}

func TranslateArithmetic(command VMCommand) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (cw *CodeWriter) WritePushPop(ctype VMCommandType, seg string, idx int) error {
	if ctype != C_PUSH && ctype != C_POP {
		return fmt.Errorf("invalid command type %d", ctype)
	}

	// translate the command
	asmcommand, err := TranslatePushPop(ctype, seg, idx, cw.fileName)
	if err != nil {
		return err
	}
	io.WriteString(cw, asmcommand)
	return nil
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
