package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type FileTreeNode struct {
	Value    os.FileInfo
	Children []*FileTreeNode
}

func NewNode(path string) (*FileTreeNode, error) {
	nodeI, err := os.Stat(path)

	node := &FileTreeNode{
		Value:    nodeI,
		Children: []*FileTreeNode{},
	}

	return node, err
}

func populateTree(root *FileTreeNode) error {
	if !root.Value.IsDir() {
		log.Println("the root is not a directory")
		return nil
	}
	entries, err := os.ReadDir(root.Value.Name())
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", root.Value.Name(), err)
	}
	for _, entry := range entries {
		childPath := filepath.Join(root.Value.Name(), entry.Name())
		childNode, err := NewNode(childPath)
		if err != nil {
			return err
		}
		root.Children = append(root.Children, childNode)
	}

	return nil
}

func walk(curr *FileTreeNode, path *[]os.FileInfo) *[]os.FileInfo {
	if curr == nil {
		fmt.Printf("no children left")
		return path
	}
	//in order traversal presumes that you will go as far left from the current head node as possible to then print it when there is nothing further

	// pre
	*path = append(*path, curr.Value)
	// recurse
	for i := 0; i < len(curr.Children); i++ {
		walk(curr.Children[i], path)
	}
	return path
}

func main() {
	rootNode, err := NewNode(".")
	if err != nil {
		log.Fatalf("there is an error reading stat of root node tree: %v\n", err)
	}
	if err := populateTree(rootNode); err != nil {
		fmt.Printf("Error populating children: %v\n", err)
		return
	}
	// Assuming you have a root node of the file tree
	var filePaths []os.FileInfo
	path := walk(rootNode, &filePaths)
	log.Println(path)
	// Now filePaths contains the paths of all files in the tree
	for _, filePath := range filePaths {
		fmt.Println(": ", filePath.Name())
	}
}
