package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
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
	}
)

// error checking:
func check(e error) {
	if e != nil {
		log.Fatalf("Some error occured: %v\n", e)
	}
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

func NewFileNode(value os.FileInfo, relPath string) *FileTreeNode {
	return &FileTreeNode{Value: value, Path: relPath}
}

func (n *FileTreeNode) BFS() *FileTreeNode {
	//..procedure BFS(G, root) is
	//2      let Q be a queue
	//3      label root as explored
	//4      Q.enqueue(root)
	//5      while Q is not empty do
	//6          v := Q.dequeue()
	//7          if v is the goal then
	//8              return v
	//9          for all edges from v to w in G.adjacentEdges(v) do
	//10              if w is not labeled as explored then
	//11                  label w as explored
	//12                  w.parent := v
	//13                  Q.enqueue(w)
	q := NewQueue[FileTreeNode]()
	q.Enqueue(*n)

	for q.Length != 0 {
		v := q.Dequeue()
		if v.changed() {
			return v
		}
		if v.Value.IsDir() {
			files, err := os.ReadDir(v.Path)
			check(err)
			for i := 0; i < len(files); i++ {
				//check later and change
				f, err := files[i].Info()
				check(err)
				var path string
				if files[i].IsDir() {
					path = "1"
				} else {
					path = "2"
				}
				q.Enqueue(*NewFileNode(f, path))
			}
		} else {
			log.Fatalln("The root is not a directory, try changing your root")
		}
	}

	return n
}

func (v *FileTreeNode) changed() bool {
	return false
}

func (n *FileTreeNode) listFiles() {
	if n.Value.IsDir() {
		files, err := os.ReadDir(n.Path)
		check(err)
		for _, f := range files {
			t, err := f.Info()
			check(err)
			node := NewFileNode(t, filepath.Join(n.Value.Name(), f.Name()))
			n.Children = append(n.Children, node)
		}
	} else {

	}
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
}
