package shell

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type File struct {
	PID         int
	ProcessName string
}

func newFile(processName string) *File {
	return &File{ProcessName: processName}
}

func CreateFile() (*File, error) {
	nf := newFile(fmt.Sprintf("godemon_%s", RandChar()))
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	gDirPath := fmt.Sprintf("%s/.godemon/", root)
	// Check if the folder exists
	if _, err := os.Stat(gDirPath); os.IsNotExist(err) {
		// If the folder does not exist, create it
		err := os.MkdirAll(gDirPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return nil, err
		}
		fmt.Println("Directory created successfully")
	}
	cmd := exec.Command(fmt.Sprintf("go build -o %s%s", gDirPath, nf.ProcessName))
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	pids, err := GetPIDs(nf.ProcessName)
	if err != nil {
		return nil, err
	}
	fmt.Println(pids)
	nf.PID = pids[0]
	return nf, nil
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

//
//// // LINUX ONLY (maybe Mac, haven't tried).
//func RestartL[T any]() {
//	qPID := new(q.Queue[T])
//
//}
//func RestartW() {
//
//}
