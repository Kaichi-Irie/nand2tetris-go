package main

import "nand2tetris-golang/ch06/hack"

func main() {
	err := hack.Hack()
	if err != nil {
		panic(err)
	}
}
