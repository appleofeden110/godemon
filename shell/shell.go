package shell

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func Shell() error {
	rootPath, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("Shell error, getting absolute path: %v\n", err)
	}
	prgName, err := os.Stat(rootPath)
	if err != nil {
		return fmt.Errorf("Shell error, getting root path: %v\n", err)
	}
	buildCmd := fmt.Sprintf("go build -o %s %s", prgName.Name(), rootPath)

	pid, err := getPidByName(prgName.Name())
	if err != nil {
		return fmt.Errorf("Shell error, error checking if the app is still running: %v\n", err)
	}

	if pid > 0 {
		if err := stopProcess(pid); err != nil {
			fmt.Errorf("Shell error, failed to stop the running instance: %v\n", err)
		}
	}

	if err := execCommand(buildCmd); err != nil {
		return fmt.Errorf("Shell error, build failed")
	}

	startCmd := fmt.Sprintf("./%s", prgName.Name())
	if err := execCommand(startCmd); err != nil {
		return fmt.Errorf("Shell error, failed to start the application: %v\n", err)
	}
	log.Println("Process restarted")
	fmt.Println("Process restarted")
	return nil
}

func execCommand(cmnd string) error {
	cmd := exec.Command("bash", "-c", cmnd)
	out := new(bytes.Buffer)
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %s", err, out.String())
	}
	fmt.Println(out.String())
	return nil
}
func getPidByName(processName string) (int, error) {
	// Execute the tasklist command to find the process by name
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH", "/FI", fmt.Sprintf("IMAGENAME eq %s", processName))
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to run tasklist: %v", err)
	}

	// Parse the output
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		// CSV output is used for easier parsing
		fields := strings.Split(line, "\",\"")
		if len(fields) > 1 {
			// Remove leading and trailing " characters
			fields[0] = strings.TrimPrefix(fields[0], "\"")
			fields[1] = strings.Trim(fields[1], "\"\r")
			if strings.ToLower(fields[0]) == strings.ToLower(processName) {
				// Found the process, return its PID
				pid, err := strconv.Atoi(fields[1])
				if err != nil {
					return 0, fmt.Errorf("failed to parse PID '%s' as integer: %v", fields[1], err)
				}
				return pid, nil
			}
		}
	}

	// Process not found
	return 0, nil
}

func stopProcess(pid int) error {
	// Convert the PID to a string
	pidStr := strconv.Itoa(pid)

	// Prepare the taskkill command
	cmd := exec.Command("taskkill", "/PID", pidStr, "/F")
	out := new(bytes.Buffer)
	cmd.Stdout = out
	cmd.Stderr = out
	// Execute the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to terminate process with PID %s: %v", pidStr, err)
	}

	return nil
}
