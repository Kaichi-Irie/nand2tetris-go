package jackanalyzer

import (
	"fmt"
	"nand2tetris-go/jackcompiler/compilationengine"
	"os"
	"path/filepath"
)

func Analize(path string) error {
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
		xmlFilePath := jackFilePath[:len(jackFilePath)-5] + ".xml"
		jackFile, err := os.Open(jackFilePath)
		if err != nil {
			return err
		}
		defer jackFile.Close()

		xmlFile, err := os.Create(xmlFilePath)
		if err != nil {
			return err
		}
		defer xmlFile.Close()
		ce := compilationengine.CreateCEwithFirstToken(xmlFile, jackFile)
		err = ce.CompileClass()
		if err != nil {
			return err
		}
		fmt.Println("compiled", jackFilePath, "to", xmlFilePath)
	}
	fmt.Println("done")
	return nil
}
