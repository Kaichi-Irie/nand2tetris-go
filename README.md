![img](img/nand2tetris-go.jpg)
![Coverage](https://img.shields.io/badge/Coverage-66.0%25-yellow)
# Nand2TetrisをGolangで実装する
[O'Reilly Japan - コンピュータシステムの理論と実装 第2版](https://www.oreilly.co.jp/books/9784814400874/)（通称Nand2Tetris）の実装プロジェクトを勉強がてら，go言語で実装していきます．パッケージのドキュメントは[こちら](https://kaichi-irie.github.io/nand2tetris-go/pkg/nand2tetris-go/index.html)．書評は[こちら](https://qiita.com/garudakai/items/7e09c95ef8b2a3c4e8be)．

※ 1-5章はOnline IDEで実装
- [x] 1章: ブール論理
- [x] 2章: ブール算術
- [x] 3章: メモリー
- [x] 4章: 機械語
- [x] 5章: コンピュータアーキテクチャ
- [x] 6章: アセンブラ
- [x] 7章: VM（コンパイラ バックエンド）1: 処理
- [x] 8章: VM（コンパイラ バックエンド）2: 制御
- [x] 9章: 高水準言語
- [x] 10章: コンパイラ（コンパイラ フロントエンド）1: 構文解析
- [x] 11章: コンパイラ（コンパイラ フロントエンド）2: コード生成
- [ ] 12章: OS

## Nand2Tetrisとは
　Nand2Tetrisプロジェクトとは，Nandゲートから始めて，論理ゲート、加算器、CPUを設計したのち、アセンブラ、VM変換器、コンパイラ、OSを実装し，コンピュータを完成させます．そして，最後にその上でアプリケーション（テトリスなど）を動作させるというプロジェクトです．

　本プロジェクトでは，Hackという専用のコンピュータアーキテクチャ，およびその上で動作するJackという専用の高水準言語が用意されています．また，前半のハードウェアプロジェクトは，公式に[Online IDE](https://nand2tetris.github.io/web-ide)が用意されており，そこで実装することができます．後半（6章以降）のソフトウェアプロジェクトは，各自，好きなプログラミング言語で実装することができます．

　プロジェクトの概要については，公式のCourse Promoの動画がとてもわかりやすいです．
[![Course Promo](https://img.youtube.com/vi/wTl5wRDT0CU/0.jpg)](https://youtu.be/wTl5wRDT0CU?si=cpyPA9cG7uHAp2tA "Course Promo")

## プロジェクトの構成
- `assembler/`: Hackアセンブラの実装
- `vm/`: VM変換器の実装
- `compiler/`: Jackコンパイラの実装
- `img/`: プロジェクトの画像
- `docs/`: 自動生成されたドキュメント
- `generate-doc.sh`: ドキュメントを自動生成するスクリプト
# 実行方法

## テスト
テストは，`go test`コマンドを使って実行します．
```sh
$ go test ./...
```


## アセンブラ
アセンブラは，Hackアセンブリ言語をHackバイナリファイルに変換するプログラムです．
アセンブリファイル（.asm）を引数に与えて実行すると，同じディレクトリに`<input.asm>.hack`が生成されます．これはHackバイナリファイルです．
```sh
$ cd assembler
$ go run main.go <input.asm>
```

## VM変換器
VM変換器は，Hack VM言語をHackアセンブリ言語に変換するプログラムです．
VMファイル（.vm）を引数に与えて実行すると，同じディレクトリに`<input.vm>.asm`が生成されます．これはHackアセンブリファイルです．

```sh
$ cd vm
$ go run main.go <input.vm>
```

ディレクトリを引数に与えると，ディレクトリ内の全てのVMファイルを変換し，同じディレクトリに`<dirname>.asm`が生成されます．この場合には，`Sys.vm`が最初に実行されるアセンブリコードが生成されることに注意してください．`Sys.vm`が含まれていない場合は，エラーが発生します．
```sh
$ go run main.go <dirname>
```

## Jackコンパイラ
Jackコンパイラは，Jack言語をHack VM言語に変換するプログラムです．コンパイラは構文解析とコード生成の2つのフェーズに分かれています．
構文解析では，Jack言語のソースコードを解析し，構文木を生成します．その後，構文木をHack VM言語に変換します．
```sh
$ cd compiler
$ go run main.go <input.jack>
```
このXMLファイルは，`<input>.vm`という名前で出力されます．
ディレクトリを引数に与えると，ディレクトリ内の全てのJackファイルをコンパイルし，同じディレクトリに`<filename>.vm`が生成されます．
```sh
$ go run main.go <dirname>
```
## ドキュメント生成
godocのinstall
```
$ go install golang.org/x/tools/cmd/godoc@latest
$ export GOPATH=$HOME/go
$ export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
godocの実行
```
$./generate-doc.sh
```

# References
- [nand2tetris](https://www.nand2tetris.org/)
- [O'Reilly Japan - コンピュータシステムの理論と実装 第2版](https://www.oreilly.co.jp/books/9784814400874/)
- [The Elements of Computing Systems: Building a Modern Computer from First Principles: Nisan, Noam, Schocken, Shimon](https://www.amazon.com/Elements-Computing-Systems-Building-Principles/dp/0262640686)
- [GitHub - YadaYuki/nand2tetris-golang: Nand2tetris Implementation by Goʕ◔ϖ◔ʔ 😺](https://github.com/YadaYuki/nand2tetris-golang)
    - 同じくgo言語で実装されたNand2Tetrisのプロジェクト．
