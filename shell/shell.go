package shell

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
)

type File struct {
	PID         int32
	ProcessName string
}

func GetPIDs(processName string) []int {
	cmd := exec.Command("pgrep", fmt.Sprintf("%v", processName))

	r, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("There is an error invoking Stdout pipe: %v\n", err)
		return nil
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)

	if err := cmd.Start(); err != nil {
		log.Fatalf("There is an error starting a command: %v\n", err)
		return nil
	}
	pids := make([]int, 0)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		i++
		pid, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("there is an error converting: %v\n", err)
			return nil
		}
		pids = append(pids, pid)
		fmt.Printf("%v: %v\n", i, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("There is an error with the scanner: %v\n", err)
		return nil
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("Error running the command: %v\n", err)
		return nil
	}
	return pids
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
