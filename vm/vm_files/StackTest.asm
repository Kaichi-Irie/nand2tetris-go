// push constant 17
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 17
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
// eq
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@EQ_0_TRUE
D;JEQ
(EQ_0_FALSE)
D=0
@EQ_0_END
0;JMP
(EQ_0_TRUE)
D=-1
@EQ_0_END
0;JMP
(EQ_0_END)
@SP
A=M
M=D
@SP
M=M+1
// push constant 17
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 16
@16
D=A
@SP
A=M
M=D
@SP
M=M+1
// eq
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@EQ_1_TRUE
D;JEQ
(EQ_1_FALSE)
D=0
@EQ_1_END
0;JMP
(EQ_1_TRUE)
D=-1
@EQ_1_END
0;JMP
(EQ_1_END)
@SP
A=M
M=D
@SP
M=M+1
// push constant 16
@16
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 17
@17
D=A
@SP
A=M
M=D
@SP
M=M+1
// eq
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@EQ_2_TRUE
D;JEQ
(EQ_2_FALSE)
D=0
@EQ_2_END
0;JMP
(EQ_2_TRUE)
D=-1
@EQ_2_END
0;JMP
(EQ_2_END)
@SP
A=M
M=D
@SP
M=M+1
// push constant 892
@892
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 891
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
// lt
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@LT_3_TRUE
D;JLT
(LT_3_FALSE)
D=0
@LT_3_END
0;JMP
(LT_3_TRUE)
D=-1
@LT_3_END
0;JMP
(LT_3_END)
@SP
A=M
M=D
@SP
M=M+1
// push constant 891
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 892
@892
D=A
@SP
A=M
M=D
@SP
M=M+1
// lt
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@LT_4_TRUE
D;JLT
(LT_4_FALSE)
D=0
@LT_4_END
0;JMP
(LT_4_TRUE)
D=-1
@LT_4_END
0;JMP
(LT_4_END)
@SP
A=M
M=D
@SP
M=M+1
// push constant 891
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 891
@891
D=A
@SP
A=M
M=D
@SP
M=M+1
// lt
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@LT_5_TRUE
D;JLT
(LT_5_FALSE)
D=0
@LT_5_END
0;JMP
(LT_5_TRUE)
D=-1
@LT_5_END
0;JMP
(LT_5_END)
@SP
A=M
M=D
@SP
M=M+1
// push constant 32767
@32767
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 32766
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
// gt
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@GT_6_TRUE
D;JGT
(GT_6_FALSE)
D=0
@GT_6_END
0;JMP
(GT_6_TRUE)
D=-1
@GT_6_END
0;JMP
(GT_6_END)
@SP
A=M
M=D
@SP
M=M+1
// push constant 32766
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 32767
@32767
D=A
@SP
A=M
M=D
@SP
M=M+1
// gt
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@GT_7_TRUE
D;JGT
(GT_7_FALSE)
D=0
@GT_7_END
0;JMP
(GT_7_TRUE)
D=-1
@GT_7_END
0;JMP
(GT_7_END)
@SP
A=M
M=D
@SP
M=M+1
// push constant 32766
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 32766
@32766
D=A
@SP
A=M
M=D
@SP
M=M+1
// gt
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@GT_8_TRUE
D;JGT
(GT_8_FALSE)
D=0
@GT_8_END
0;JMP
(GT_8_TRUE)
D=-1
@GT_8_END
0;JMP
(GT_8_END)
@SP
A=M
M=D
@SP
M=M+1
// push constant 57
@57
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 31
@31
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 53
@53
D=A
@SP
A=M
M=D
@SP
M=M+1
// add
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D+M
@SP
A=M
M=D
@SP
M=M+1
// push constant 112
@112
D=A
@SP
A=M
M=D
@SP
M=M+1
// sub
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D-M
@SP
A=M
M=D
@SP
M=M+1
// neg
@SP
M=M-1
A=M
D=M
D=-D
@SP
A=M
M=D
@SP
M=M+1
// and
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D&M
@SP
A=M
M=D
@SP
M=M+1
// push constant 82
@82
D=A
@SP
A=M
M=D
@SP
M=M+1
// or
@SP
M=M-1
A=M
D=M
@R13
M=D
@SP
M=M-1
A=M
D=M
@R13
D=D|M
@SP
A=M
M=D
@SP
M=M+1
// not
@SP
M=M-1
A=M
D=M
D=!D
@SP
A=M
M=D
@SP
M=M+1
// infinite loop
(END)
@END
0;JMP
