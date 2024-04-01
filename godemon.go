package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
			fls[KeyFile{childNode.Value.Name(), childPath}] = childNode.Value.ModTime()
			_, _ = BLR(childPath, fls)

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

func mapCompare(map1, map2 map[KeyFile]time.Time) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, val := range map1 {
		if val2, exists := map2[key]; !exists || !val.Equal(val2) {
			return false
		}
	}
	return true
}

func GodemonInit() error {
	flsBackUp := make(map[KeyFile]time.Time)
	fls := make(map[KeyFile]time.Time)

	for {
		_, err := BLR(".", fls)
		check(err)

		if !mapCompare(flsBackUp, fls) {
			err := shell()
			if err != nil {
				fmt.Printf("There is an error running bash file: %v\n", err)
				return err
			}
			// Create a new map and deep copy fls into it
			newBackup := make(map[KeyFile]time.Time)
			for k, v := range fls {
				newBackup[k] = v
			}
			flsBackUp = newBackup // Now flsBackUp is a deep copy of fls

			log.Println("Files have changed, action taken.")
		} else {
			log.Println("No changes detected.")
		}

		fmt.Println("Waiting for changes...")
		time.Sleep(400 * time.Millisecond)
	}
}

//
//func Main() {
//	err := GodemonInit()
//	if err != nil {
//		log.Fatalf("pizdec: %v\n", err)
//	}
//}

func main() {
	err := GodemonInit()
	if err != nil {
		panic(err)
	}
}
