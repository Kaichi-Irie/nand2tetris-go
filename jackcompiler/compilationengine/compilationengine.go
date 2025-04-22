package compilationengine

import (
	"io"
	st "nand2tetris-go/jackcompiler/symboltable"
	tk "nand2tetris-go/jackcompiler/tokenizer"
)

type CompilationEngine struct {
	writer       io.Writer
	t            *tk.Tokenizer
	classST      *st.SymbolTable
	subroutineST *st.SymbolTable
}

func New(xmlFile io.Writer, r io.Reader) *CompilationEngine {
	return &CompilationEngine{
		writer:       xmlFile,
		t:            tk.New(r),
		classST:      st.NewSymbolTable(),
		subroutineST: st.NewSymbolTable(),
	}
}

func NewWithFirstToken(xmlFile io.Writer, r io.Reader) *CompilationEngine {
	t, err := tk.NewWithFirstToken(r)
	if err != nil {
		panic(err)
	}
	return &CompilationEngine{
		writer:       xmlFile,
		t:            t,
		classST:      st.NewSymbolTable(),
		subroutineST: st.NewSymbolTable(),
	}
}
