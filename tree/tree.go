package tree

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/appleofeden110/godemon/utils"
	"io"
	"os"
	"path/filepath"
	"time"
)

type (
	FileTreeNode struct {
		Value    os.FileInfo
		Path     string
		Children []*FileTreeNode
	}

	KeyFile struct {
		Name, Path string
	}
)

var (
	ErrChanged = errors.New("File has changed")
)

func (n *FileTreeNode) Error(err error) error {
	return fmt.Errorf("GoDemon 𓁹‿𓁹: %v", err)
}

func NewFileNode(relPath string) *FileTreeNode {
	value, err := os.Stat(relPath)
	utils.Check(err)
	return &FileTreeNode{Value: value, Path: relPath}
}

// checks the file tree and gives FileTreeNode with the whole tree in it, and otherwise gives an error
func BLR(path string, fls map[KeyFile]time.Time) (*FileTreeNode, error) {
	n := NewFileNode(path)
	if n.Value.IsDir() {
		files, err := os.ReadDir(path)
		utils.Check(err)
		dirignore := make(map[string]bool)
		err = IgnoreDirs(dirignore)
		utils.Check(err)
		for _, f := range files {
			if f.IsDir() && dirignore[f.Name()] {
				continue
			}
			childPath := filepath.Join(path, f.Name())
			childNode := NewFileNode(childPath)
			fls[KeyFile{childNode.Value.Name(), childPath}] = childNode.Value.ModTime()
			_, _ = BLR(childPath, fls)

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
		return fmt.Errorf("error: %v", err)
	}
	err = json.Unmarshal(b, &ignoreDirs)
	if err != nil {
		return fmt.Errorf("There is an error unmarshaling data: %v\n", err)
	}
	return nil
}
