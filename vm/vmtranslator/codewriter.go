package vmtranslator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CodeWriter translates VM commands to Hack assembly code and writes the code to an output file.
type CodeWriter struct {
	File         io.Writer
	VmFileStem   string         // the base name of the .vm file without the .vm extension. e.g. "SimpleAdd"
	CommandCount int            // for generating unique labels
	FunctionName string         // for generating unique return labels
	ReturnCount  map[string]int // for generating unique return labels
}

// NewCodeWriter creates a new asm file with the given path and returns a CodeWriter. CodeWriter.FileNameStem is set to "", so it must be set before calling WriteCommand.
func NewCodeWriter(asmFile io.Writer) *CodeWriter {
	return &CodeWriter{
		File:        asmFile,
		ReturnCount: make(map[string]int),
	}
}

func NewCodeWriterFromFile(asmFilePath string) (*CodeWriter, error) {
	if filepath.Ext(asmFilePath) != ".asm" {
		return nil, fmt.Errorf("invalid file extension")
	}
	asmFile, err := os.Create(asmFilePath)
	if err != nil {
		return nil, err
	}
	return NewCodeWriter(asmFile), nil
}

// Write writes the given bytes to the output file.
func (cw *CodeWriter) Write(b []byte) (int, error) {
	return cw.File.Write(b)
}

// push the value in D register to the stack; RAM[SP]=D, SP++
var push_D = "@SP\nA=M\nM=D\n@SP\nM=M+1\n"

// pop the value from the stack to D register; SP--, D=RAM[SP]
var pop_D = "@SP\nM=M-1\nA=M\nD=M\n"

// pop the value from the stack to RAM[13]; SP--, RAM[13]=RAM[SP]
var pop_R13 = pop_D + "@R13\nM=D\n"

// TranslatePushPop generates the assembly code for VMcommand "push segment index" or "pop segment index". fileName is the name of the .vm file. It returns an error if the segment is invalid.
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

// TranslateArithmetic generates the assembly code for VMcommand "add", "sub", "and", "or", "neg", "not", "eq", "gt", or "lt". cnt is a counter for generating unique labels and used for eq, gt, and lt commands.
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

// TranslateLabel generates the assembly code for VMcommand "label label". label is the label name.
func TranslateLabel(label string) (string, error) {
	return fmt.Sprintf("(%s)\n", label), nil
}

// TranslateGoto generates the assembly code for VMcommand "goto label". label is the label to jump to.
func TranslateGoto(label string) (string, error) {
	return fmt.Sprintf("@%s\n0;JMP\n", label), nil
}

// TranslateIf generates the assembly code for VMcommand "if-goto label". label is the label to jump to if the top of the stack is not zero.
func TranslateIf(label string) (string, error) {
	asmcommand := pop_D
	asmcommand += fmt.Sprintf("@%s\nD;JNE\n", label)
	return asmcommand, nil
}

// TranslateFunction generates the assembly code for VMcommand "function functionName nVars". functionName is the name of the function, and nVars is the number of local variables.
func TranslateFunction(functionName string, nVars int) (string, error) {
	if nVars < 0 {
		return "", fmt.Errorf("nVars must be non-negative")
	}
	asmcommand, err := TranslateLabel(functionName)
	if err != nil {
		return "", err
	}
	for i := range nVars {
		command1, err := TranslatePushPop(C_PUSH, "constant", 0, "")
		if err != nil {
			return "", err
		}
		command2, err := TranslatePushPop(C_POP, "local", i, "")
		if err != nil {
			return "", err
		}
		command3, err := TranslatePushPop(C_PUSH, "local", i, "")
		if err != nil {
			return "", err
		}
		asmcommand += command1 + command2 + command3
	}
	return asmcommand, nil
}

// TranslateCall generates the assembly code for VMcommand "call functionName nArgs". functionName is the name of the function, nArgs is the number of arguments, and cnt is a counter for generating unique return labels.
func TranslateCall(functionName string, nArgs int, cnt int) (string, error) {
	// push return address
	returnAddress := fmt.Sprintf("%s$ret.%d", functionName, cnt)
	asmcommand := fmt.Sprintf("@%s\nD=A\n", returnAddress)
	asmcommand += push_D
	// push LCL, ARG, THIS, THAT
	for _, seg := range []string{"LCL", "ARG", "THIS", "THAT"} {
		asmcommand += fmt.Sprintf("@%s\nD=M\n", seg)
		asmcommand += push_D
	}
	// ARG=SP-nArgs-5
	asmcommand += fmt.Sprintf("@5\nD=A\n@%d\nD=D+A\n@SP\nD=M-D\n@ARG\nM=D\n", nArgs)
	// LCL=SP
	asmcommand += "@SP\nD=M\n@LCL\nM=D\n"
	// goto functionName
	gotoCommand, err := TranslateGoto(functionName)
	if err != nil {
		return "", err
	}
	asmcommand += gotoCommand
	// (returnAddress)
	labelCommand, err := TranslateLabel(returnAddress)
	if err != nil {
		return "", err
	}
	asmcommand += labelCommand
	return asmcommand, nil
}

// TranslateReturn generates the assembly code for VMcommand "return".
func TranslateReturn() (string, error) {
	asmcommand := "@LCL\nD=M\n@R14\nM=D\n"                 // R14=FRAME=LCL
	asmcommand += "@R14\nD=M\n@5\nA=D-A\nD=M\n@R15\nM=D\n" // R15=RET=*(FRAME-5)
	asmcommand += pop_D                                    // pop return value to D
	asmcommand += "@ARG\nA=M\nM=D\n"                       // *ARG=pop()
	asmcommand += "@ARG\nD=M+1\n@SP\nM=D\n"                // SP=ARG+1
	// restore THAT, THIS, ARG, LCL
	asmcommand += "@R14\nD=M\n@1\nA=D-A\nD=M\n@THAT\nM=D\n"
	asmcommand += "@R14\nD=M\n@2\nA=D-A\nD=M\n@THIS\nM=D\n"
	asmcommand += "@R14\nD=M\n@3\nA=D-A\nD=M\n@ARG\nM=D\n"
	asmcommand += "@R14\nD=M\n@4\nA=D-A\nD=M\n@LCL\nM=D\n"
	asmcommand += "@R15\nA=M\n0;JMP\n" // goto RET
	return asmcommand, nil
}

// resolveLabel resolves the label name for a function and a label base. If functionName is empty, it returns labelBase. Otherwise, it returns functionName$labelBase.
func resolveLabel(functionName string, labelBase string) string {
	if functionName == "" {
		return labelBase
	}
	return functionName + "$" + labelBase
}

// WriteCommand writes the assembly code for the given VM command to the output file. It returns an error if the command is invalid. It also updates the internal state of the CodeWriter, which is used for generating unique labels.
func (cw *CodeWriter) WriteCommand(gotoCommand VMCommand) error {
	// output the command as a comment
	io.WriteString(cw, "// "+string(gotoCommand)+"\n")
	ctype := getCommandType(gotoCommand)
	var asmcommand string
	var err error
	switch ctype {
	case C_ARITHMETIC:
		asmcommand, err = TranslateArithmetic(gotoCommand, cw.CommandCount)
		cw.CommandCount++
	case C_PUSH, C_POP:
		if cw.VmFileStem == "" {
			return fmt.Errorf("fileNameStem is not set")
		}
		asmcommand, err = TranslatePushPop(ctype, arg1(gotoCommand), arg2(gotoCommand), cw.VmFileStem)
	case C_LABEL:
		label := resolveLabel(cw.FunctionName, arg1(gotoCommand))
		asmcommand, err = TranslateLabel(label)
	case C_GOTO:
		label := resolveLabel(cw.FunctionName, arg1(gotoCommand))
		asmcommand, err = TranslateGoto(label)
	case C_IF:
		label := resolveLabel(cw.FunctionName, arg1(gotoCommand))
		asmcommand, err = TranslateIf(label)
	case C_FUNCTION:
		asmcommand, err = TranslateFunction(arg1(gotoCommand), arg2(gotoCommand))
		cw.FunctionName = arg1(gotoCommand)
	case C_CALL:
		cnt, ok := cw.ReturnCount[arg1(gotoCommand)]
		if !ok {
			cnt = 0
			cw.ReturnCount[arg1(gotoCommand)] = 0
		}
		asmcommand, err = TranslateCall(arg1(gotoCommand), arg2(gotoCommand), cnt)
		cw.ReturnCount[arg1(gotoCommand)]++
	case C_RETURN:
		asmcommand, err = TranslateReturn()
	default:
		return fmt.Errorf("invalid command type %d", ctype)
	}
	if err != nil {
		return err
	}
	_, err = io.WriteString(cw, asmcommand)
	return err
}

// WriteInfinityLoop writes an infinite loop to the output file. It is used to prevent the program from exiting.
func (cw *CodeWriter) WriteInfinityLoop() error {
	// TODO: Avoid label name collision
	label := "INFINITE_LOOP_END"
	_, err := io.WriteString(cw, fmt.Sprintf("// infinite loop\n(%s)\n@%s\n0;JMP\n", label, label))
	return err
}

// WriteBootStrap writes the bootstrap code to the output file. It initializes the stack pointer and calls Sys.init.
func (cw *CodeWriter) WriteBootStrap() error {
	_, err := io.WriteString(cw, "// bootstrap code\n")
	if err != nil {
		return err
	}
	// SP=256
	_, err = io.WriteString(cw, "@256\nD=A\n@SP\nM=D\n")
	if err != nil {
		return err
	}
	cw.VmFileStem = "Sys"
	return cw.WriteCommand("call Sys.init 0")
}
