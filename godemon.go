package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type (
	Qnode[T any] struct {
		Value T
		Next  *Qnode[T]
		Prev  *Qnode[T]
	}
	Queue[T any] struct {
		Head   *Qnode[T]
		Tail   *Qnode[T]
		Length int8
	}
	Qinterface[T any] interface {
		Enqueue(v T)
	}
	FileTreeNode struct {
		Value    os.FileInfo
		Path     string
		Children []*FileTreeNode
	}
	FileChecking[T any] interface {
		BFS(v os.FileInfo)
		Error() error
	}
)

func (n *FileTreeNode) Error() error {
	return fmt.Errorf("There is an error in godemon 𓁹‿𓁹")
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Enqueue(v T) {
	newNode := Qnode[T]{Value: v}
	if q.Head == nil {
		q.Head = &newNode
		q.Tail = &newNode
	} else {
		q.Tail = &newNode
		q.Tail.Next = &newNode
	}
	q.Length++
}

func (q *Queue[T]) Dequeue() *T {
	if q.Head == nil {
		return nil
	}

	head := q.Head
	q.Head = q.Head.Next

	head.Next = nil
	q.Length--
	return &head.Value
}

func NewFileNode(relPath string) *FileTreeNode {
	value, err := os.Stat(relPath)
	check(err)
	return &FileTreeNode{Value: value, Path: relPath}
}

func (n *FileTreeNode) BFS(t time.Time) *FileTreeNode {

	q := NewQueue[FileTreeNode]()
	q.Enqueue(*n)

	for q.Length != 0 {
		v := q.Dequeue()
		if v.changed(t) {
			return v
		}
		if v.Value.IsDir() {
			files, err := os.ReadDir(v.Path)
			check(err)
			for i := 0; i < len(files); i++ {
				//check later and change
				var path string
				if files[i].IsDir() {
					path = "1"
				} else {
					path = "2"
				}
				q.Enqueue(*NewFileNode(path))
			}
		} else {
			log.Fatalln("The root is not a directory, try changing your root")
		}
	}

	return n
}

func (v *FileTreeNode) changed(t time.Time) bool {
	if v.Value.ModTime() != t {
		return true
	}
	return false
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

func (root *FileTreeNode) GodemonInit() error {
	for {
		initTime := time.Now()
		root.BFS(initTime)
		// logic for restart and ongoing
		if root.Error() != nil {
			fmt.Printf("There is an error in programme: %v\n", root.Error())
			return root.Error()
		}
	}
}

func main() {
	n := NewFileNode(".")
	err := n.GodemonInit()
	check(err)
}
