package godemon

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
	// Now, proceed with the build
	absPath, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("1.1: %v", err)
	}
	rootName, err := os.Stat(absPath)
	if err != nil {
		fmt.Println("there is a problem calling os.stat: ", err)
		return err
	}

	cmd := exec.Command("go", "build", ".")
	cmd.Dir = absPath // Set the working directory if needed
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to build: %v", err)
	}
	fmt.Println("process built")

	cmd = exec.Command("cmd", "/C", rootName.Name())
	if err := cmd.Run(); err != nil { // Use Start instead of Run
		return fmt.Errorf("Failed to execute exe: %v", err)
	}
	cmd = exec.Command("cmd", "exit 0")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed to exit cmd: %v\n", err)
	}
	fmt.Println("Server Restarted")

	return nil
}

func (n *FileTreeNode) Error(err error) error {
	return fmt.Errorf("GoDemon 𓁹‿𓁹: %v", err)
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
