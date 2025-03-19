// push constant 111
@111
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 333
@333
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 888
@888
D=A
@SP
A=M
M=D
@SP
M=M+1
// pop static 8
@SP
M=M-1
A=M
D=M
@./vm/vm_files/StaticTest.8
M=D
// pop static 3
@SP
M=M-1
A=M
D=M
@./vm/vm_files/StaticTest.3
M=D
// pop static 1
@SP
M=M-1
A=M
D=M
@./vm/vm_files/StaticTest.1
M=D
// push static 3
@./vm/vm_files/StaticTest.3
D=M
@SP
A=M
M=D
@SP
M=M+1
// push static 1
@./vm/vm_files/StaticTest.1
D=M
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
// push static 8
@./vm/vm_files/StaticTest.8
D=M
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
// infinite loop
(END)
@END
0;JMP
