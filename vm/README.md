# プロジェクト: VMTranslator

## Parser
Description: `parser.go` は、VM言語のコマンドをパースするための構造体と関数を提供する。

### Types:

| types           | description                                                                                                                                                                                         |
| --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `VMCommand`     | コマンドを表す文字列。例: `"push constant 10"`, `"add"`, `"label LOOP"`, `"goto LOOP"`, `"if-goto LOOP"`                                                                                            |
| `VMCommandType` | コマンドのタイプを表す列挙型。`C_ARITHMETIC`, `C_PUSH`, `C_POP`, `C_LABEL`, `C_GOTO`, `C_IF`, `C_FUNCTION`, `C_RETURN`, `C_CALL`                                                                    |
| `CodeScanner`   | コードのコマンドをスキャンするための構造体．一行ごとにコマンドを読み込む．コメント行と空白行は無視する．また，コマンドの先頭と末尾の空白文字を取り除き，複数の空白文字を1つの空白文字に置き換える． |
| `Parser`        | VM言語のコマンドをパースするための構造体。                                                                                                                                                          |

### Functions and methods:

| functions/methods                                               | arguments   | return values   | description                                                                                                                                                                                     |
| --------------------------------------------------------------- | ----------- | --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `(p *Parser) advance`                                           |             | `bool`          | `Parser` が次の命令を読み込み、それを現在の命令にする。もしもう命令がない場合は false を返す。 コメント行や空行は無視して，命令が見つかる，またはファイルの終わりに達するまでスキャンを続ける。 |
| `getcommandType`                                                | `VMCommand` | `VMCommandType` | 命令のタイプを返す。                                                                                                                                                                            |
| `arg1`                                                          | `VMCommand` | `string`        | VMコマンドを受け取り，そのコマンドの最初の引数を返す。`C_ARITHMETIC` の場合は命令自体を返す。                                                                                                   |
| `arg2`                                                          | `VMCommand` | `int`           | VMコマンドを受け取り，そのコマンドの2番目の引数を返す。コマンドタイプが`C_PUSH`, `C_POP`, `C_FUNCTION`, `C_CALL` の場合のみ有効。                                                               |
| `NewParser(r io.Reader, commentPrefix string) Parser`           |             |                 | 新しい `Parser` 構造体を作成する。                                                                                                                                                              |
| `NewCodeScanner(r io.Reader, commentPrefix string) CodeScanner` |             |                 | 新しい `CodeScanner` 構造体を作成する。                                                                                                                                                         |

<!-- TODO: CodeWriterのドキュメントを完了する -->
## CodeWriter
Description: `codewriter.go` は、VM言語のコマンドをHackアセンブリ言語に変換するための構造体と関数を提供する。

### Types:
| types       | description                                                                                   |
| ----------- | --------------------------------------------------------------------------------------------- |
| `CodeWriter` | Hackアセンブリ言語に変換するための構造体。`io.Writer` インターフェースを実装している。 |


### Functions and methods:
| functions/methods                            | arguments                     | return values    | description                                                                                                                |
| -------------------------------------------- | ----------------------------- | ---------------- | -------------------------------------------------------------------------------------------------------------------------- |
| `TranslatePushPop`                           | `string`, `int`,  `string`    | `string`,`error` | PushまたはPopコマンドをHackアセンブリ言語に変換する。`fileName` 現在のVMファイル名を受け取るのは，staticセグメントのため． |
| `TranslateArithmetic`                        | `VMCommand`                   | `string`         | 算術コマンドをHackアセンブリ言語に変換する。                                                                               |
| `(cw *CodeWriter) WriteCommand`              | `VMCommand`                   | `error`          | VMコマンドをHackアセンブリ言語に変換し、ファイルに書き込む。 アセンブリ言語には，コメントとしてVMコマンドも追加される。    |
| `(cw *CodeWriter) WritePushPop`              | `VMCommandType`, `seg`, `idx` | `error`          | PushまたはPopコマンドをHackアセンブリ言語に変換し、ファイルに書き込む。                                                    |
| `(cw *CodeWriter) WriteArithmetic`           | `VMCommand`                   | `error`          | 算術コマンドをHackアセンブリ言語に変換し、ファイルに書き込む。                                                             |
| `NewCodeWriter(fileName string) *CodeWriter` |                               |                  | 新しい `CodeWriter` 構造体を作成する。                                                                                     |
