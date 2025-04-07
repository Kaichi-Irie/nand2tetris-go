package compilationengine

import (
	"io"
	"nand2tetris-go/jackcompiler/tokenizer"
)

type CompilationEngine struct {
	xmlFile io.Writer
	t       *tokenizer.Tokenizer
}

func New(xmlFile io.Writer, r io.Reader) *CompilationEngine {
	return &CompilationEngine{
		xmlFile: xmlFile,
		t:       tokenizer.New(r),
	}
}

func CreateCEwithFirstToken(xmlFile io.Writer, r io.Reader) *CompilationEngine {
	t, err := tokenizer.CreateTokenizerWithFirstToken(r)
	if err != nil {
		panic(err)
	}
	return &CompilationEngine{
		xmlFile: xmlFile,
		t:       t,
	}
}
