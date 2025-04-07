package tokenizer

var XMLEscapes = map[string]string{
	LESS.Val():    "&lt;",
	GREATER.Val(): "&gt;",
	AND.Val():     "&amp;",
	"\"":          "&quot;",
}
