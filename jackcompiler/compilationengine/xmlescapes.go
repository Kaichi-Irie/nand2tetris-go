package compilationengine

import (
	tk "nand2tetris-go/jackcompiler/tokenizer"
)

var XMLEscapes = map[string]string{
	tk.LESS.Val:    "&lt;",
	tk.GREATER.Val: "&gt;",
	tk.AND.Val:     "&amp;",
	"\"":           "&quot;",
}
