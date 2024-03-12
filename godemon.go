package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type KeyFile struct {
	Name, Path string
}

func BLR(path string, fls map[KeyFile]time.Time) (*FileTreeNode, error) {
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
			fmt.Printf("name: %v\n path: %v\n modification time: %v\n", childNode.Value.Name(), childPath, childNode.Value.ModTime())
			fls[KeyFile{childNode.Value.Name(), childPath}] = childNode.Value.ModTime()
			_, _ = BLR(childPath, fls)

			if childNode != nil {
				n.Children = append(n.Children, childNode)
			}
		}
	}
	fmt.Printf("Hello, buddy:%v\n", fls[KeyFile{"dir2", "Dir/dir2"}])
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

//
//func mapCompare(map1, map2 map[KeyFile]time.Time) bool {
//	if len(map1) != len(map2) {
//		return false
//	}
//	for key, val := range map1 {
//		if val2, exists := map2[key]; !exists || !val.Equal(val2) {
//			return false
//		}
//	}
//	return true
//}
//
//func GodemonInit() error {
//	flsBackUp := make(map[KeyFile]time.Time)
//	fls := make(map[KeyFile]time.Time)
//	for {
//		root, err := BLR(".", fls)
//		if len(flsBackUp) == 0 || mapCompare(flsBackUp, fls) {
//
//		}
//		if checkf("There is some trouble with app", err) {
//			log.Fatalln("")
//			return err
//		}
//		fmt.Println("Hui NAZ")
//		time.Sleep(400 * time.Millisecond)
//	}
//}

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

func main() {
	fls := make(map[KeyFile]time.Time)

	thing, err := BLR(".", fls)
	check(err)
	fmt.Println(thing.Path, thing.Value.ModTime())

	for i := 0; i < len(thing.Children); i++ {
		fmt.Println(thing.Children[i].Value.Name())
	}
	//err = shell()
	//check(err)
}
