package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileTreeNode struct {
	Value    os.FileInfo
	Children []*FileTreeNode
}

func NewFileTree(paths []string, path string, ignoreDirs map[string]bool) (*FileTreeNode, []string, error) {
	abs, err := filepath.Abs(path)
	log.Println(abs)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v\n", err)
		return nil, nil, err
	}
	node, err := os.Stat(abs)
	if err != nil {
		log.Fatalf("Error getting stat of a file: %v\n", err)
		return nil, nil, err
	}
	//init the node of the tree
	n := &FileTreeNode{
		Value:    node,
		Children: make([]*FileTreeNode, 0),
	}

	if !node.IsDir() {
		return n, paths, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("There is a problem in reading the directory: %v\n", err)
		return nil, nil, err
	}
	log.Println(entries)
	for i := 0; i < len(entries); i++ {
		if ignoreDirs[entries[i].Name()] && entries[i].IsDir() {
			continue
		}
		np := filepath.Join(abs, entries[i].Name())
		paths = append(paths, np)
		log.Println("Paths: ", paths)
		path = np
		g, ps, err := NewFileTree(paths, path, ignoreDirs)
		if err != nil {
			log.Fatalf("error creating the tree inside recursion: %v\n", err)
			return nil, nil, err
		}
		log.Println(ps)
		n.Children = append(n.Children, g)
	}
	return n, paths, nil
}

func walk(curr *FileTreeNode, path *[]string) []string {
	if curr == nil {
		fmt.Printf("no children left")
		return *path
	}
	//in order traversal presumes that you will go as far left from the current head node as possible to then print it when there is nothing further

	// pre
	*path = append(*path, curr.Value.Name())
	// recurse
	for i := 0; i < len(curr.Children); i++ {
		walk(curr.Children[i], path)
	}
	return *path
}

func IgnoreDirs(ignoreDirs map[string]bool) error {
	jsonF, err := os.Open("ignoreDIrs.json")
	if err != nil {
		log.Fatalf("There is an error reading json file: %v\n", err)
		return err
	}
	defer jsonF.Close()

	b, err := io.ReadAll(jsonF)
	if err != nil {
		log.Fatalf("error: %v\n", err)
		return err
	}
	err = json.Unmarshal(b, &ignoreDirs)
	if err != nil {
		log.Fatalf("There is an error unmarshaling data: %v\n", err)
		return err
	}
	return nil
}

func main() {
	d := make(map[string]bool, 0)
	err := IgnoreDirs(d)
	if err != nil {
		log.Fatalln("THere is an error: ", err)
		return
	}
	log.Println(d)
	ps := make([]string, 0)
	thing, paths, err := NewFileTree(ps, ".", d)
	if err != nil {
		log.Fatalf("Fuck: ", err)
		return
	}
	log.Printf("thing: %v\n, thing2: %v\n", thing.Value.Name(), paths)
}
