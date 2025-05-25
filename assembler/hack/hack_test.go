package hack

import (
	"bytes"
	"strings"
	"testing"
)

func TestHack(t *testing.T) {
	tests := []struct {
		asmCode  string
		expected string
	}{
		{asmCode: `@1
D=-1

// comment

@2
AM=D
// comment2
@0
0;JMP`, expected: `0000000000000001
1110111010010000
0000000000000010
1110001100101000
0000000000000000
1110101010000111
`},
		{asmCode: `// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/6/add/Add.asm

// Computes R0 = 2 + 3  (R0 refers to RAM[0])

@2
D=A
@3
D=D+A
@0
M=D
`, expected: `0000000000000010
1110110000010000
0000000000000011
1110000010010000
0000000000000000
1110001100001000
`},
		{asmCode: `// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/6/rect/Rect.asm

// Draws a rectangle at the top-left corner of the screen.
// The rectangle is 16 pixels wide and R0 pixels high.
// Usage: Before executing, put a value in R0.

   // If (R0 <= 0) goto END else n = R0
   @R0
   D=M
   @END
   D;JLE
   @n
   M=D
   // addr = base address of first screen row
   @SCREEN
   D=A
   @addr
   M=D
(LOOP)
   // RAM[addr] = -1
   @addr
   A=M
   M=-1
   // addr = base address of next screen row
   @addr
   D=M
   @32
   D=D+A
   @addr
   M=D
   // decrements n and loops
   @n
   MD=M-1
   @LOOP
   D;JGT
(END)
   @END
   0;JMP
`, expected: `0000000000000000
1111110000010000
0000000000010111
1110001100000110
0000000000010000
1110001100001000
0100000000000000
1110110000010000
0000000000010001
1110001100001000
0000000000010001
1111110000100000
1110111010001000
0000000000010001
1111110000010000
0000000000100000
1110000010010000
0000000000010001
1110001100001000
0000000000010000
1111110010011000
0000000000001010
1110001100000001
0000000000010111
1110101010000111
`},
	}

	for _, test := range tests {
		t.Run(test.asmCode, func(t *testing.T) {
			testHack(t, test.asmCode, test.expected)
		})
	}
}
func testHack(t *testing.T, asmCode, expected string) {
	// Create a reader for the asmCode
	asmFile := strings.NewReader(asmCode)
	hackFile := &bytes.Buffer{}
	err := Hack(asmFile, hackFile)
	if err != nil {
		t.Fatalf("Hack failed: %v", err)
	}

	if hackFile.String() != expected {
		t.Errorf("Expected %s, got %s", expected, hackFile.String())
	}
}
