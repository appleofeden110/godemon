package tree_check

import (
	"fmt"
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

func (n *FileTreeNode) Error(err error) error {
	return fmt.Errorf("Godemon error: %v\n", err)
}

func NewFileNode(relPath string) (*FileTreeNode, error) {
	value, err := os.Stat(relPath)
	if err != nil {
		return nil, err
	}
	return &FileTreeNode{Value: value, Path: relPath}, nil
}

// checks the file tree and gives FileTreeNode with the whole tree in it, and otherwise gives an error
func BLR(path string, fls map[KeyFile]time.Time) (*FileTreeNode, error) {
	n, err := NewFileNode(path)

	if err != nil {
		return nil, err
	}
	if n.Value.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		dirignore := make(map[string]bool)
		err = ignoreDirs(dirignore)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if f.IsDir() && dirignore[f.Name()] {
				continue
			}
			childPath := filepath.Join(path, f.Name())
			childNode, err := NewFileNode(childPath)
			if err != nil {
				return nil, err
			}
			fls[KeyFile{childNode.Value.Name(), childPath}] = childNode.Value.ModTime()
			_, _ = BLR(childPath, fls)

			if childNode != nil {
				n.Children = append(n.Children, childNode)
			}
		}
	}
	return n, nil
}
