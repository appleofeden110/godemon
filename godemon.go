package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func BLR(path string) (*FileTreeNode, error) {
	n := newFileNode(path)
	if n.Value.IsDir() {
		files, err := os.ReadDir(path)
		check(err)
		dirignore := make(map[string]bool, 0)
		err = IgnoreDirs(dirignore)
		fmt.Println()
		check(err)
		for _, f := range files {
			if f.IsDir() && dirignore[f.Name()] {
				continue
			}
			childPath := filepath.Join(path, f.Name())
			childNode := newFileNode(childPath)
			fmt.Printf("name: %v\n path: %v\n", f.Name(), childNode.Path)

			_, _ = BLR(childPath)
			if childNode != nil {
				n.Children = append(n.Children, childNode)
			}
		}
	}
	return n, nil
}

func IgnoreDirs(ignoreDirs map[string]bool) error {
	jsonF, err := os.Open("ignoreDirs.json")
	if err != nil {
		return fmt.Errorf("There is an error reading json file: %v\n", err)

	}
	defer jsonF.Close()

	b, err := io.ReadAll(jsonF)
	if err != nil {
		return fmt.Errorf("error", err)
	}
	err = json.Unmarshal(b, &ignoreDirs)
	if err != nil {
		return fmt.Errorf("There is an error unmarshaling data: %v\n", err)
	}
	return nil
}

func main() {
	thing, err := BLR(".")
	check(err)
	fmt.Println(thing.Path, thing.Value.ModTime())

	for i := 0; i < len(thing.Children); i++ {
		fmt.Println(thing.Children[i].Value.Name())
	}

}
