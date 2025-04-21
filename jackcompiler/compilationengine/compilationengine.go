package compilationengine

import (
	"io"
	tk "nand2tetris-go/jackcompiler/tokenizer"
)

type CompilationEngine struct {
	writer io.Writer
	t      *tk.Tokenizer
}

func New(xmlFile io.Writer, r io.Reader) *CompilationEngine {
	return &CompilationEngine{
		writer: xmlFile,
		t:      tk.New(r),
	}
}

func NewWithFirstToken(xmlFile io.Writer, r io.Reader) *CompilationEngine {
	t, err := tk.NewWithFirstToken(r)
	if err != nil {
		panic(err)
	}
	return &CompilationEngine{
		writer: xmlFile,
		t:      t,
	}
}
