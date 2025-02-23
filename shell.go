package godemon

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	// "github.com/appleofeden110/godemon/queue"
)

//go:embed embed/godemon_cli
var cliBinary []byte

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
	var pids []int

	// Read /proc directory
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("error reading /proc: %v", err)
	}

	for _, entry := range entries {
		// Skip if not a number (PIDs are numbers)
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		// Read the process name from /proc/[pid]/comm
		commPath := fmt.Sprintf("/proc/%d/comm", pid)
		data, err := os.ReadFile(commPath)
		if err != nil {
			continue
		}

		// Compare process name (trim newline)
		name := strings.TrimSpace(string(data))
		if name == processName {
			pids = append(pids, pid)
		}
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

// // StartDetachedProcess should have an argument of the current program that is running this function, so that godemon can
// // simply
func StartDetachedProcess(args []string) error {
	pid := os.Getpid()
	err := Godemon_log_pid(pid, "Starting GoDemon process")
	if err != nil {
		return fmt.Errorf("Error for godemon logging: %v\n", err)
	}
	// to do: same terminal output, daemon logic
	return nil
}

// Godemon_log to show that it is a logging of godemon, this function makes sure it: writes GODEMON: at the start and a newline after the message
func Godemon_log(processName string, thingsToWrite ...string) error {
	pids, err := GetPIDs(processName)
	if err != nil {
		return fmt.Errorf("error getting pid for log: %v\n", err)
	}
	pid := pids[0]

	processStdIn := fmt.Sprintf("/proc/%d/fd/1", pid)
	for i := 0; i < len(thingsToWrite); i++ {
		err = os.WriteFile(processStdIn, []byte("GODEMON: "+thingsToWrite[i]+"\n"), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error writing file with process name of %s, error: %v\n", processName, err)
		}
	}
	return nil
}

// Godemon_log to show that it is a logging of godemon, this function makes sure it: writes GODEMON: at the start and a newline after the message
func Godemon_log_pid(pid int, thingsToWrite ...string) error {
	processStdIn := fmt.Sprintf("/proc/%d/fd/1", pid)
	for i := 0; i < len(thingsToWrite); i++ {
		err := os.WriteFile(processStdIn, []byte("GODEMON: "+thingsToWrite[i]+"\n"), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error writing file with pid of %d, error: %v\n", pid, err)
		}
	}
	return nil
}
