package jackanalyzer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Kaichi-Irie/nand2tetris-go/jackcompiler/compilationengine"
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
		className := jackFilePath[:len(jackFilePath)-5]
		jackFile, err := os.Open(jackFilePath)
		if err != nil {
			return err
		}
		defer jackFile.Close()
		vmFile, err := os.Create(className + ".vm")
		if err != nil {
			return err
		}

		ce := compilationengine.NewWithVMWriter(vmFile, jackFile, className)
		err = ce.CompileClass()
		if err != nil {
			return err
		}
		fmt.Println("compiled", jackFilePath, "to", vmFile.Name())
	}
	fmt.Println("done")
	return nil
}
