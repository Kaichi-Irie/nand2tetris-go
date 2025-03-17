package hack

import (
	"testing"
)

func TestSymbol(t *testing.T) {
	st := NewSymbolTable()
	constant := SymbolOrConstant("100")
	code, err := symbol(constant, &st)
	if code != "000000001100100" || err != nil {
		t.Fatalf(`symbol("100") = %q, %v, want "000000001100100", nil`, code, err)
	}

	st.table["LABEL"] = 8
	sym := SymbolOrConstant("LABEL")
	code, err = symbol(sym, &st)
	if code != "000000000001000" || err != nil {
		t.Fatalf(`symbol("LABEL") = %q, %v, want "000000000001000", nil`, code, err)
	}
}

func TestDest(t *testing.T) {
	mnemonic := Mnemonic("M")
	code, err := dest(mnemonic)
	if code != "001" || err != nil {
		t.Fatalf(`dest("M") = %q, %v, want "001", nil`, code, err)
	}

	mnemonic = "D"
	code, err = dest(mnemonic)
	if code != "010" || err != nil {
		t.Fatalf(`dest("D") = %q, %v, want "010", nil`, code, err)
	}
}

func TestComp(t *testing.T) {
	mnemonic := Mnemonic("D")
	code, err := comp(mnemonic)
	if code != "0001100" || err != nil {
		t.Fatalf(`comp("D") = %q, %v, want "0001100", nil`, code, err)
	}
}

func TestJump(t *testing.T) {
	mnemonic := Mnemonic("JGT")
	code, err := jump(mnemonic)
	if code != "001" || err != nil {
		t.Fatalf(`jump("JGT") = %q, %v, want "001", nil`, code, err)
	}
}
