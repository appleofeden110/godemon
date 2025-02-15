package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	// "github.com/appleofeden110/godemon/queue"
)

type File struct {
	PID         int
	ProcessName string
	path        string
}

func newFile(processName string) *File {
	return &File{ProcessName: processName}
}

// CreateFile creates the new build of the godemon file with randomized characted to it.
// It is using the newFile function that is creating new instance of a File struct.
func CreateFile() (*File, error) {
	nf := newFile(fmt.Sprintf("godemon_%s", RandChar()))
	// I think it is better to have a tmp folder for the whole system to be able to process it basically. (it will delete it on the end of the run of the program)
	gDirPath := "/tmp/.godemon/"

	// Check if the folder exists
	if _, err := os.Stat(gDirPath); os.IsNotExist(err) {
		// If the folder does not exist, create it
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
		err := os.MkdirAll(gDirPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return nil, errors.New(fmt.Sprintf("Error with creating the .godemon: %v\n", err))
		}
		fmt.Println("Directory created successfully")
	}
	fullPath := fmt.Sprintf("%s%s", gDirPath, nf.ProcessName)
	cmdCommand := []string{"build", "-o", fullPath}
	cmd := exec.Command("go", cmdCommand...)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	pids, err := GetPIDs(nf.ProcessName)
	if err != nil {
		return nil, err
	}
	nf.PID = pids[0]
	nf.path = fullPath
	return nf, nil
}

func (f *File) RunProc() error {
	cmd := exec.Command(fmt.Sprintf("%v", f.path))
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Problem with cmd starting the process: %v\n", err)
	}
	return nil
}

func (f *File) SuspendProc() error {
	cmd := exec.Command("kill", strconv.Itoa(f.PID))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running the kill command: %v\n", err)
	}
	return nil
}

func GetPIDs(processName string) ([]int, error) {
	cmd := exec.Command("pgrep", fmt.Sprintf("%v", processName))

	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("There is an error invoking Stdout pipe: %v\n", err))
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)

	if err := cmd.Start(); err != nil {
		return nil, errors.New(fmt.Sprintf("There is an error starting a command: %v\n", err))
	}
	pids := make([]int, 0)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		i++
		pid, err := strconv.Atoi(line)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("there is an error converting: %v\n", err))
		}
		pids = append(pids, pid)
		fmt.Printf("%v: %v\n", i, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("There is an error with the scanner: %v\n", err))
	}

	if err := cmd.Wait(); err != nil {
		return nil, errors.New(fmt.Sprintf("Error running the command: %v\n", err))
	}
	return pids, nil
}

// Returns 6 random characters, used to name the files differently.
// Counted just for fun: probability of two different files being named the same at the end of the day is:
// 0.0000000026268% or 2.6268*10^-9%.
// Basically very very little chance of the "merge error". Also LINUX ONLY (for now)
func RandChar() string {
	chars := make([]byte, 0)
	for i := 0; i < 6; i++ {
		r := rand.Intn(58)
		c := byte(r + 65)
		chars = append(chars, c)
	}
	return string(chars)
}

// StartDetachedProcess should have an argument of the current program that is running this function, so that godemon can
// simply
func StartDetachedProcess(args []string) error {
	pids, err := GetPIDs(args[0])
	if err != nil {
		return fmt.Errorf("error starting detached process, pids step: %v\n", err)
	}
	if len(pids) > 1 {
		log.Println("found more than one process, using the one with pid: ", pids[0])
	}
	// to do: same terminal output, daemon logic
	return nil
}
