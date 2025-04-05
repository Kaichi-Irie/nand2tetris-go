package jackanalyzer

import (
	"fmt"
	"nand2tetris-golang/jackcompiler/tokenizer"
	"os"
	"path/filepath"
)

func Analize(path string) error {
	// path can be a .vm file or a directory containing .vm files
	fmt.Println("jack analyzer")
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	var jackFilePaths []string

	if info.IsDir() {
		jackFilePaths, err = filepath.Glob(filepath.Join(path, "*.jack"))
		if err != nil {
			return err
		}
	} else if filepath.Ext(path) == ".jack" {
		jackFilePaths = []string{path}
	} else {
		return fmt.Errorf("input path must be a .jack file or a directory")
	}
	for _, jackFilePath := range jackFilePaths {
		vmFilePath := jackFilePath[:len(jackFilePath)-5] + ".vm"
		jackFile, err := os.Open(jackFilePath)
		if err != nil {
			return err
		}
		defer jackFile.Close()
		// TODO: implement the analyzer
		tokenizer := tokenizer.NewTokenizer(jackFile)
		for tokenizer.advance() {
			token := tokenizer.currentToken
			tokenType, err := tokenizer.GetTokenType(token)
			if err != nil {
				return err
			}
			// TODO: implement the token type handling
		}

		vmFile, err := os.Create(vmFilePath)
		if err != nil {
			return err
		}
		defer vmFile.Close()

	}
	fmt.Println("done")
	return nil
}
