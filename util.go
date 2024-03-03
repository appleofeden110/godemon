package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type (
	FileTreeNode struct {
		Value    os.FileInfo
		Path     string
		Children []*FileTreeNode
	}
	FileChecking[T any] interface {
		BFS(t time.Time) *FileTreeNode
		Error() error
	}
)

func (n *FileTreeNode) Error() error {
	return fmt.Errorf("There is an error in godemon 𓁹‿𓁹")
}

func newFileNode(relPath string) *FileTreeNode {
	value, err := os.Stat(relPath)
	check(err)
	return &FileTreeNode{Value: value, Path: relPath}
}

// check() takes in error and checks if it's nil or not.
// It takes the boilerplate and basically transforms into a very easy function.
// If needed, error can always be specified, no need to call it on every error in the program
func check(e error) {
	if e != nil {
		log.Fatalf("Some error occured: %v\n", e)
	}
}

func checkf(msg string, e error) {
	if e != nil {
		log.Fatalf("%v: %v\n", msg, e)
	}
}

// changed(t time.Time) (bool) tracks actual changing time of the file.
// "t" represents initial time that everything tracks with.
func (v *FileTreeNode) changed(t time.Time) bool {
	if v.Value.ModTime() != t {
		return true
	}
	return false
}
