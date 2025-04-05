package compilationengine

import (
	"bytes"
	"nand2tetris-go/jackcompiler/tokenizer"
	"strings"
	"testing"
)

func TestCompileClass(t *testing.T) {
	xmlFile := &bytes.Buffer{}
	ce := New(xmlFile, strings.NewReader(""))

	tests := []struct {
		jackCode    string
		expectedXML string
	}{
		{
			jackCode: `class Main {static int i;static int j;}`,
			expectedXML: `<class>
<keyword> class </keyword>
<identifier> Main </identifier>
<symbol> { </symbol>
<keyword> static </keyword>
<keyword> int </keyword>
<identifier> i </identifier>
<symbol> ; </symbol>
<keyword> static </keyword>
<keyword> int </keyword>
<identifier> j </identifier>
<symbol> ; </symbol>
<symbol> } </symbol>
</class>
`,
		},
	}
	var err error
	for _, test := range tests {
		ce.t, err = tokenizer.CreateTokenizerWithFirstToken(strings.NewReader(test.jackCode))
		if err != nil {
			t.Errorf("CreateTokenizerWithFirstToken() error: %v", err)
		}
		err = ce.CompileClass()
		if err != nil {
			t.Errorf("CompileClass() error: %v", err)
		}
		// remove leading and trailing whitespace from the actual XML
		if xmlFile.String() != test.expectedXML {
			t.Errorf("CompileClass() = %v, want %v", xmlFile.String(), test.expectedXML)
		}
		xmlFile.Reset()
	}

}
