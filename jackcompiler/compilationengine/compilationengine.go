package compilationengine

import (
	"io"

	st "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/symboltable"
	tk "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/tokenizer"
	vw "github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/vmwriter"
)

type CompilationEngine struct {
	vmwriter     *vw.VMWriter
	t            *tk.Tokenizer
	classST      *st.SymbolTable
	subroutineST *st.SymbolTable
	labelCount   int // for generating unique labels
}

func New(vmwriter io.Writer, r io.Reader, className string) *CompilationEngine {
	return &CompilationEngine{
		vmwriter:     vw.New(vmwriter),
		t:            tk.New(r),
		classST:      st.NewSymbolTable(),
		subroutineST: st.NewSymbolTable(),
		labelCount:   0,
	}
}

func NewWithFirstToken(vmwriter io.Writer, r io.Reader, className string) *CompilationEngine {
	ce := New(vmwriter, r, className)
	t, err := tk.NewWithFirstToken(r)
	if err != nil {
		panic(err)
	}
	ce.t = t
	return ce
}

func NewWithVMWriter(vmwriter io.Writer, r io.Reader, className string) *CompilationEngine {
	ce := NewWithFirstToken(vmwriter, r, className)
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
