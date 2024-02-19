package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var paths = make([]string, 0)

type FileTreeNode struct {
	Value    os.FileInfo
	Children []*FileTreeNode
}

func NewTree(path string, ignoreDirs map[string]bool) (*FileTreeNode, error) {
	nodeI, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	node := &FileTreeNode{
		Value:    nodeI,
		Children: []*FileTreeNode{},
	}

	if !nodeI.IsDir() {
		// If the current path is not a directory, return the node without children
		return node, nil
	}

	entries, err := os.ReadDir(path)
	paths = append(paths, path)

	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Check if the directory should be ignored
			if ignoreDirs[entry.Name()] {
				continue // Skip this directory
			}
		}
		childPath := filepath.Join(path, entry.Name())
		childNode, err := NewTree(childPath, ignoreDirs)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, childNode)
	}

	return node, nil
}

func walk(curr *FileTreeNode, path *[]os.FileInfo) []os.FileInfo {
	if curr == nil {
		fmt.Printf("no children left")
		return *path
	}
	//in order traversal presumes that you will go as far left from the current head node as possible to then print it when there is nothing further

	// pre
	*path = append(*path, curr.Value)
	// recurse
	for i := 0; i < len(curr.Children); i++ {
		walk(curr.Children[i], path)
	}
	return *path
}

func main() {
	iD := map[string]bool{
		".git":  true,
		".idea": true,
	}
	rootNode, err := NewTree(".", iD)
	log.Println(paths)
	for _, i := range rootNode.Children {
		fmt.Printf("hui %v\n", i.Value.IsDir())
		fmt.Printf("hui2 %v\n", i.Value.Name())
		if i.Value.IsDir() {
			files, err := os.ReadDir(i.Value.Name())
			if err != nil {
				log.Fatalln("error reading directory: ", err)
			}
			fmt.Printf("something: %v\n", files)

		}
	}

	if err != nil {
		log.Fatalf("there is an error reading stat of root node tree: %v\n", err)
	}
	r, err := os.Stat(".")
	var filePaths []os.FileInfo
	root := &FileTreeNode{
		Value:    r,
		Children: []*FileTreeNode{},
	}
	// Assuming you have a root node of the file tree
	p := walk(root, &filePaths)
	for _, i := range p {
		println(i.Name())
	}
}
