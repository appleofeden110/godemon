package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func BLR(path string) (*FileTreeNode, error) {
	n := newFileNode(path)
	if n.Value.IsDir() {
		files, err := os.ReadDir(path)
		check(err)
		dirignore := make(map[string]bool, 0)
		err = IgnoreDirs(dirignore)
		check(err)
		for _, f := range files {
			if f.IsDir() && dirignore[f.Name()] {
				continue
			}
			childPath := filepath.Join(path, f.Name())
			childNode := newFileNode(childPath)
			fmt.Printf("name: %v\n path: %v\n modTime: %v\n", f.Name(), childNode.Path, childNode.Value.ModTime())

			_, _ = BLR(childPath)
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
		return fmt.Errorf("error", err)
	}
	err = json.Unmarshal(b, &ignoreDirs)
	if err != nil {
		return fmt.Errorf("There is an error unmarshaling data: %v\n", err)
	}
	return nil
}

//func (root *FileTreeNode) GodemonInit() error {
//  for {
//    initTime := time.Now()
//    root.BFS(initTime)
//    // logic for restart and ongoing
//    if root.Error() != nil {
//      fmt.Printf("There is an error in programme: %v\n", root.Error())
//      return root.Error()
//    }
//  }
//}

func shell() error {
	absPath, err := filepath.Abs(".")
	check(err)
	rootName, err := os.Stat(absPath)
	check(err)
	scriptName := "restart.bat"
	// Define the contents of the batch file
	contents := fmt.Sprintf(`
		echo hello world
		go build %v
		%v\%v.exe
	`, absPath, absPath, rootName)

	// Create the batch file
	file, err := os.Create("restart.bat")
	if err != nil {
		log.Println("there is an error: ", err)
		return err
	}
	defer file.Close()

	// Write the contents to the file
	_, err = file.WriteString(contents)
	if err != nil {
		panic(err)
	}

	// Execute the batch file
	err = exec.Command(filepath.Join(absPath, scriptName)).Run()
	if err != nil {
		fmt.Printf("Failed to execute restart.bat: %v\n", err)
		return err
	}
	fmt.Println("Server Restarted")
	return nil
}

func main() {
	thing, err := BLR(".")
	check(err)
	fmt.Println(thing.Path, thing.Value.ModTime())

	for i := 0; i < len(thing.Children); i++ {
		fmt.Println(thing.Children[i].Value.Name())
	}

}
