package compilationengine

import (
	"bytes"
	"io"
	st "nand2tetris-go/jackcompiler/symboltable"
	tk "nand2tetris-go/jackcompiler/tokenizer"
	vw "nand2tetris-go/jackcompiler/vmwriter"
)

type CompilationEngine struct {
	writer       io.Writer
	vmwriter     *vw.VMWriter
	t            *tk.Tokenizer
	classST      *st.SymbolTable
	subroutineST *st.SymbolTable
	labelCount   int // for generating unique labels
}

func New(xmlFile io.Writer, r io.Reader, className string) *CompilationEngine {
	return &CompilationEngine{
		writer:       xmlFile,
		vmwriter:     vw.New(bytes.NewBuffer(nil)),
		t:            tk.New(r),
		classST:      st.NewSymbolTable(),
		subroutineST: st.NewSymbolTable(),
		labelCount:   0,
	}
}

func NewWithFirstToken(xmlFile io.Writer, r io.Reader, className string) *CompilationEngine {
	ce := New(xmlFile, r, className)
	t, err := tk.NewWithFirstToken(r)
	if err != nil {
		panic(err)
	}
	ce.t = t
	return ce
}

func NewWithVMWriter(vmwriter io.Writer, xmlFile io.Writer, r io.Reader, className string) *CompilationEngine {
	ce := NewWithFirstToken(xmlFile, r, className)
	ce.vmwriter = vw.New(vmwriter)
	return ce
}

func (ce *CompilationEngine) Lookup(name string) (st.Identifier, bool) {
	if id, ok := ce.subroutineST.Lookup(name); ok {
		return id, true
	}
	if id, ok := ce.classST.Lookup(name); ok {
		return id, true
	}
	return st.Identifier{}, false

}
