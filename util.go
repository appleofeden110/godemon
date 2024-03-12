package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type (
	FileTreeNode struct {
		Value    os.FileInfo
		Path     string
		Children []*FileTreeNode
	}
)

var (
	ErrChanged = errors.New("File has changed")
)

func shell() error {
	absPath, err := filepath.Abs(".")
	checkf("1.1", err)
	rootName, err := os.Stat(absPath)
	checkf("1.2", err)
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
	checkf("hui", err)

	// Execute the batch file
	err = exec.Command(filepath.Join(absPath, scriptName)).Run()
	if err != nil {
		fmt.Printf("Failed to execute restart.bat: %v\n", err)
		return err
	}
	fmt.Println("Server Restarted")
	return nil
}

func (n *FileTreeNode) Error(err error) error {
	return fmt.Errorf("There is an error in godemon 𓁹‿𓁹: %v", err)
}

func newFileNode(relPath string) *FileTreeNode {
	value, err := os.Stat(relPath)
	check(err)
	return &FileTreeNode{Value: value, Path: relPath}
}

// check() takes in error and checks if it's nil or not.
// It takes the boilerplate and basically transforms into a very easy function.
// If needed, error can always be specified, no need to call it on every error in the program
func check(e error) bool {
	if e != nil {
		log.Fatalf("Some error occured: %v\n", e)
		return true
	}
	return false
}

func checkf(msg string, e error) bool {
	if e != nil {
		log.Fatalf("%v: %v\n", msg, e)
		return true
	}
	return false
}

// changed(t time.Time) (bool) tracks actual changing time of the file.
// "t" represents initial time that everything tracks with.
//func (n *FileTreeNode) changed() bool {
//	if n.Value.ModTime() !=  {
//		return true
//	}
//	return false
//}
