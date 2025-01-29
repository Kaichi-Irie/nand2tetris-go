
# プロジェクト: Hackアセンブラ

Description:
`hack` パッケージは、Hackコンピュータのアセンブリ言語プログラムをバイナリコードに変換するための関数を提供する。

## `hack.go`

| functions | arguments | return values | description |
|-----------|-----------|---------------|-------------|
| `Hack` | `string` |`error`| アセンブリ言語のファイルをバイナリコードに変換する。 ファイル名はコマンドライン引数として渡される。 書き込みファイルは入力ファイルと同じ名前で、拡張子が`.hack`になる。 |
| `firstPass` |`string`| `SymbolTable`,`error`| 1回目のパス。L命令を探し、それらをシンボルテーブルに追加する。 |
| `secondPass` | `string`, `*os.File`, `SymbolTable` |`error`| 2回目のパス。A命令とC命令を探し、それらをバイナリコードに変換する。 |

## `symbol_table.go`
Description: `symbol_table.go` は、シンボルテーブルを表す構造体と関数を提供する。


| types | description |
|-------|-------------|
| `SymbolTable` | シンボルテーブルを表す構造体。 変数用のRAM空間を管理するカウンタ`variableCount`と、シンボルテーブルを表す`map[SymbolOrConstant]int`を持つ。 |


| functions/methods | arguments | return values | description |
|-----------|-----------|---------------|-------------|
| `NewSymbolTable` | | `SymbolTable` | 新しいシンボルテーブルを返す。 |
| `(s *SymbolTable) contains` | `SymbolOrConstant` | `bool` | シンボルテーブルに引数のシンボルが含まれているかどうかを返す。 |
| `(s *SymbolTable) addVariable` | `SymbolOrConstant` | | シンボルテーブルに変数を追加する。 |
| `(s *SymbolTable) addLabel` | `SymbolOrConstant`, `int` | | シンボルテーブルにラベルを追加する。 |

## `code.go`
Description: `code.go` は、アセンブリ言語の命令をバイナリコードに変換するための関数を提供する。



| types | description |
|-------|-------------|
| `BinaryCode` | 15ビットのバイナリコードを表す文字列。例: "000000001100100" |



| functions | arguments | return values | description |
|-----------|-----------|---------------|-------------|
| `decimalToBinary` |`string`| `BinaryCode`,`error`| `string`で表された10進数を15ビットのバイナリコードに変換する。例: "100" -> "000000001100100" |
| `symbol` | `SymbolOrConstant`, `*SymbolTable` | `BinaryCode`,`error`| シンボルまたは定数を15ビットのバイナリコードに変換する。新しい変数が見つかった場合は、シンボルテーブルに追加する。例: "100" -> "000000001100100", "LOOP" -> `SymbolTable`.table["LOOP"]->"000000000000001" |
| `dest` | `Mnemonic` | `BinaryCode`,`error`| destニーモニックのバイナリコードを返す。 |
| `comp` | `Mnemonic` | `BinaryCode`,`error`| compニーモニックのバイナリコードを返す。 |
| `jump` | `Mnemonic` | `BinaryCode`,`error`| jumpニーモニックのバイナリコードを返す。 |

## `parser.go`
Description: `parser.go` は、アセンブリ言語の命令をパースするための構造体と関数を提供する。

| types | description |
|-------|-------------|
| `Instruction` | 命令を表す文字列。例: "@2", "D=M", "(LOOP)" |
|`Mnemonic`| C命令のdest, comp, jumpのニーモニックを表す文字列。例: "D", "A", "M", "D+1", "A-1", "D-A", "D&A", "JGT" |
| `SymbolOrConstant` | シンボルまたは定数を表す文字列。例: `"LOOP"`, `"100"` ,`"x"` |
| `InstructionType` | 命令のタイプを表す列挙型。`A_Instruction`, `C_Instruction`, `L_Instruction` のいずれか。 |


| functions/methods | arguments | return values | description |
|-----------|-----------|---------------|-------------|
| `(p *Parser) advance`   |           |`bool`         | `Parser` が次の命令を読み込み、それを現在の命令にする。もしもう命令がない場合は false を返す。 コメント行や空行は無視して，命令が見つかる，またはファイルの終わりに達するまでスキャンを続ける。 |
| `(p *Parser) text`    |           |`string`       | 現在の命令を返す。 |
| `(p *Parser) symbol`    |           | `SymbolOrConstant`,`error`| 現在の命令が`"@const"`の場合は`"const"`を返す。`"@variable"`の場合は`"variable"`を返す。`"(label)"`の場合は`"label"`を返す。 A命令でもL命令でもない場合はエラーを返す。 |
| `(p *Parser) dest `     |           |`Mnemonic`,`error`| 現在のC命令のdestニーモニックを返す。 |
| `(p *Parser) comp`      |           |`Mnemonic`,`error`| 現在のC命令のcompニーモニックを返す。 |
| `(p *Parser) jump`      |           |`Mnemonic`,`error`| 現在のC命令のjumpニーモニックを返す。 |
| `isConst`   | `SymbolOrConstant` |`bool`| 引数が定数かどうかを返す。 `string` を `int` に変換できるかどうかで判断する。 |
| `getInstructionType` | `Instruction` | `InstructionType` | 命令のタイプを返す。 |
| `isEmptyLine` |`string`|`bool`| 引数が空行かどうかを返す。 |
| `isCommentLine` |`string`|`bool`| 引数がコメント行かどうかを返す。 |
