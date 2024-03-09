package godemon

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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
func check(e error) {
	if e != nil {
		log.Fatalf("Some error occured: %v\n", e)
	}
}

func checkf(msg string, e error) {
	if e != nil {
		log.Fatalf("%v: %v\n", msg, e)
	}
}

// changed(t time.Time) (bool) tracks actual changing time of the file.
// "t" represents initial time that everything tracks with.
func (v *FileTreeNode) changed(t time.Time) bool {
	if v.Value.ModTime() != t {
		return true
	}
	return false
}
