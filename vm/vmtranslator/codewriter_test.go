package vmtranslator

import (
	"testing"
)

func TestTranslatePushPop(t *testing.T) {

	tests := []struct {
		ctype    VMCommandType
		seg      string
		idx      int
		fileName string
		want     string
	}{
		{C_PUSH, "local", 3, "Test", "@LCL\nD=M\n@3\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{C_POP, "local", 3, "Test", "@LCL\nD=M\n@3\nD=D+A\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nA=M\nM=D\n"},
		{C_PUSH, "argument", 2, "Test", "@ARG\nD=M\n@2\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{C_POP, "argument", 2, "Test", "@ARG\nD=M\n@2\nD=D+A\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nA=M\nM=D\n"},
		{C_PUSH, "this", 3, "Test", "@THIS\nD=M\n@3\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{C_POP, "this", 3, "Test", "@THIS\nD=M\n@3\nD=D+A\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nA=M\nM=D\n"},
		{C_PUSH, "that", 3, "Test", "@THAT\nD=M\n@3\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{C_POP, "that", 3, "Test", "@THAT\nD=M\n@3\nD=D+A\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nA=M\nM=D\n"},
		{C_PUSH, "temp", 7, "Test", "@5\nD=A\n@7\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{C_POP, "temp", 7, "Test", "@5\nD=A\n@7\nD=D+A\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nA=M\nM=D\n"},
		{C_PUSH, "pointer", 0, "Test", "@THIS\nD=M\n@0\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{C_POP, "pointer", 0, "Test", "@THIS\nD=M\n@0\nD=D+A\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nA=M\nM=D\n"},
		{C_PUSH, "pointer", 1, "Test", "@THAT\nD=M\n@1\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{C_POP, "pointer", 1, "Test", "@THAT\nD=M\n@1\nD=D+A\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nA=M\nM=D\n"},
		{C_PUSH, "static", 3, "Test", "@Test.3\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		{C_POP, "static", 3, "Test", "@SP\nM=M-1\nA=M\nD=M\n@Test.3\nM=D\n"},
		{C_PUSH, "constant", 30, "Test", "@30\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
	}
	for _, test := range tests {
		asmcommand, err := TranslatePushPop(test.ctype, test.seg, test.idx, test.fileName)
		if err != nil {
			t.Errorf("TranslatePushPop failed: %v", err)
		}
		if asmcommand != test.want {
			var pushpop string
			switch test.ctype {
			case C_PUSH:
				pushpop = "push"
			case C_POP:
				pushpop = "pop"
			}
			t.Errorf("TranslatePushPop(%q, %q, %d, %q) = %q, want %q", pushpop, test.seg, test.idx, test.fileName, asmcommand, test.want)
		}
	}
}

func TestTranslateArithmetic(t *testing.T) {
	tests := []struct {
		command VMCommand
		want    string
	}{
		// add
		{"add", "@SP\nM=M-1\nA=M\nD=M\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nD=D+M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		// sub
		{"sub", "@SP\nM=M-1\nA=M\nD=M\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nD=D-M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		// neg
		{"neg", "@SP\nM=M-1\nA=M\nD=M\nD=-D\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		// eq
		{"eq", "@SP\nM=M-1\nA=M\nD=M\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nD=D-M\n@EQ_0_TRUE\nD;JEQ\n(EQ_0_FALSE)\nD=0\n@EQ_0_END\n0;JMP\n(EQ_0_TRUE)\nD=-1\n@EQ_0_END\n0;JMP\n(EQ_0_END)\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		// gt
		{"gt", "@SP\nM=M-1\nA=M\nD=M\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nD=D-M\n@GT_0_TRUE\nD;JGT\n(GT_0_FALSE)\nD=0\n@GT_0_END\n0;JMP\n(GT_0_TRUE)\nD=-1\n@GT_0_END\n0;JMP\n(GT_0_END)\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		// lt
		{"lt", "@SP\nM=M-1\nA=M\nD=M\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nD=D-M\n@LT_0_TRUE\nD;JLT\n(LT_0_FALSE)\nD=0\n@LT_0_END\n0;JMP\n(LT_0_TRUE)\nD=-1\n@LT_0_END\n0;JMP\n(LT_0_END)\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		// and
		{"and", "@SP\nM=M-1\nA=M\nD=M\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nD=D&M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		// or
		{"or", "@SP\nM=M-1\nA=M\nD=M\n@R13\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R13\nD=D|M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
		// not
		{"not", "@SP\nM=M-1\nA=M\nD=M\nD=!D\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"},
	}
	for _, test := range tests {
		asmcommand, err := TranslateArithmetic(test.command, 0)
		if err != nil {
			t.Errorf("TranslateArithmetic failed: %v", err)
		}
		if asmcommand != test.want {
			t.Errorf("TranslateArithmetic(%q) = %q, want %q", test.command, asmcommand, test.want)
		}
	}
}
