package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileTreeElement struct {
	entry    os.FileInfo
	children []FileTreeElement
	isLast   bool
}

func getObjects(path string, printFiles bool) ([]os.FileInfo, error) {
	objects, err := ioutil.ReadDir(path)
	if printFiles {
		return objects, err
	}
	var newObjects []os.FileInfo
	for _, obj := range objects {
		if obj.IsDir() {
			newObjects = append(newObjects, obj)
		}
	}
	return newObjects, nil
}

func getTreeStruct(path string, printFiles bool) ([]FileTreeElement, error) {
	var tree []FileTreeElement
	objects, err := getObjects(path, printFiles)
	for index, obj := range objects {
		var children []FileTreeElement
		if obj.IsDir() {
			newPath := filepath.Join(path, obj.Name())
			children, err = getTreeStruct(newPath, printFiles)
		}

		isLast := len(objects)-1 == index
		if obj.IsDir() || printFiles {
			tree = append(tree, FileTreeElement{obj, children, isLast})
		}
	}
	return tree, err
}

func getSize(el FileTreeElement) string {
	var size string
	if !el.entry.IsDir() && el.entry.Size() > 0 {
		size = fmt.Sprintf(" (%db)", el.entry.Size())
	} else if !el.entry.IsDir() {
		size = " (empty)"
	}
	return size
}

func getPrefix(oldPrefix string, isLast bool) string {
	prefix := "│\t"
	if isLast {
		prefix = "\t"
	}
	return oldPrefix + prefix
}

func drawTree(writer io.Writer, tree []FileTreeElement, prefix string) error {
	for _, elem := range tree {
		arrow := "├"
		if elem.isLast {
			arrow = "└"
		}
		line := fmt.Sprintf("%s%s───%s%s\n", prefix, arrow, elem.entry.Name(), getSize(elem))
		_, err := writer.Write([]byte(line))
		if len(elem.children) > 0 {
			err = drawTree(writer, elem.children, getPrefix(prefix, elem.isLast))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func dirTree(writer io.Writer, path string, printFiles bool) error {
	tree, err := getTreeStruct(path, printFiles)
	err = drawTree(writer, tree, "")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
