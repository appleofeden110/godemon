package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

func ignoreDirs(ignoreDirs map[string]bool) error {
	jsonF, err := os.Open("ignoreDirs.json")
	if err != nil {
		return ErrIgnoreDirs
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
