package tree

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type (
	FileTreeNode struct {
		Value    os.FileInfo
		Path     string
		Children []*FileTreeNode
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
	Check(err)
	return &FileTreeNode{Value: value, Path: relPath}
}

// check() takes in error and checks if it's nil or not.
// It takes the boilerplate and basically transforms into a very easy function.
// If needed, error can always be specified, no need to call it on every error in the program
func Check(e error) bool {
	if e != nil {
		log.Fatalf("Some error occured: %v\n", e)
		return true
	}
	return false
}

func checkf(msg string, e error) bool {
	if e != nil {
		log.Fatalf("%v: %v\n", msg, e)
		return true
	}
	return false
}
