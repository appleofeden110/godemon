package godemon

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func BLR(path string, t time.Time) (*FileTreeNode, error) {

	n := newFileNode(path)
	if n.Value.IsDir() {
		files, err := os.ReadDir(path)
		check(err)
		dirignore := make(map[string]bool)
		err = IgnoreDirs(dirignore)
		check(err)
		for _, f := range files {
			if f.IsDir() && dirignore[f.Name()] {
				continue
			}
			childPath := filepath.Join(path, f.Name())
			childNode := newFileNode(childPath)

			_, _ = BLR(childPath, childNode.Value.ModTime())
			fmt.Printf("name: %v\n path: %v\n modTime: %v\n", f.Name(), childNode.Path, childNode.Value.ModTime())
			if !childNode.Value.ModTime().Equal(t) {
				log.Printf("%v: %v", ErrChanged, childPath)
				return childNode, ErrChanged
			}
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

//func GodemonInit() error {
//	prevTime := time.Now()
//	for {
//		root, err := BLR(".", prevTime)
//		if checkf("There is some trouble with app", err) {
//			switch {
//			case err == ErrChanged:
//
//			}
//		}
//
//		time.Sleep(400 * time.Millisecond)
//	}
//}
//
//func Main() {
//	prevTime := time.Now()
//	thing, err := BLR(".", prevTime)
//	check(err)
//	fmt.Println(thing.Path, thing.Value.ModTime())
//
//	for i := 0; i < len(thing.Children); i++ {
//		fmt.Println(thing.Children[i].Value.Name())
//	}
//	err = shell()
//	check(err)
//}
