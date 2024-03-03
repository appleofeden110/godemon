package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func BLR(path string, t time.Time) (*FileTreeNode, error) {
	n := newFileNode(path)
	if n.Value.IsDir() {
		files, err := os.ReadDir(path)
		check(err)
		for _, f := range files {
			if f.IsDir() && f.Name() == ".git" || f.Name() == ".idea" {
				continue
			}
			childPath := filepath.Join(path, f.Name())
			childNode := newFileNode(childPath)
			fmt.Printf("name: %v\n path: %v\n", f.Name(), childNode.Path)

			_, _ = BLR(childPath, t)
			if childNode != nil {
				n.Children = append(n.Children, childNode)
			}
		}
	}
	return n, nil
}

func IgnoreDirs(ignoreDirs map[string]bool) error {
	jsonF, err := os.Open("ignoreDIrs.json")
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

}
