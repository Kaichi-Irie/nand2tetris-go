# プロジェクト: VMTranslator

## Parser
Description: `parser.go` は、VM言語のコマンドをパースするための構造体と関数を提供する。

Types:

| types       | description                                                                                                                      |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `Command`     | コマンドを表す文字列。例: `"push constant 10"`, `"add"`, `"label LOOP"`, `"goto LOOP"`, `"if-goto LOOP"`                         | |
| `commandType` | コマンドのタイプを表す列挙型。`C_ARITHMETIC`, `C_PUSH`, `C_POP`, `C_LABEL`, `C_GOTO`, `C_IF`, `C_FUNCTION`, `C_RETURN`, `C_CALL` |

Functions and methods:

| functions/methods     | arguments | return values | description                                                                                                                                                                                     |
| --------------------- | --------- | ------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `(p *Parser) advance` |           | `bool`        | `Parser` が次の命令を読み込み、それを現在の命令にする。もしもう命令がない場合は false を返す。 コメント行や空行は無視して，命令が見つかる，またはファイルの終わりに達するまでスキャンを続ける。 |
| `getcommandType`      | `Command` | `commandType` | 命令のタイプを返す。                                                                                                                                                                            |
| `(p *Parser) arg1`                |           | `string`      | 現在の命令の最初の引数を返す。`C_ARITHMETIC` の場合は命令自体を返す。                                                                                                                           |
| `(p *Parser) arg2`                |           | `int`         | 現在の命令の2番目の引数を返す。`C_ARITHMETIC` または `C_RETURN` の場合は 0 を返す。                                                                                                             |
