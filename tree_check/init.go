package tree_check

import (
	"fmt"
	"time"
)

func TreeCheck() error {
	flsBackUp := make(map[KeyFile]time.Time)
	fls := make(map[KeyFile]time.Time)
	var err error
	i := 0
	for {
		i++
		_, err = BLR(".", fls)
		if err != nil {
			return fmt.Errorf("Tree problem: %v\n", err)
		}
		changedFiles, deletedFiles, same := mapCompare(flsBackUp, fls)

		if !same {
			newBackup := make(map[KeyFile]time.Time)
			for k, v := range fls {
				newBackup[k] = v
			}
			flsBackUp = newBackup
			changedString, deletedString := printFileArray(changedFiles, deletedFiles)
			if i > 1 {
				if deletedString != "" {
					fmt.Printf("[%v✗%v] Files or directories deleted: \n%v\n\n", BrightRed, Reset, deletedString)
				}
				if changedString != "" {
					fmt.Printf("[%v!%v] Files and directories have changed:\n%v\n\n", BrightYellow, Reset, changedString)
				}
				fmt.Printf("SIZE OF THE PROJECT: %v\nWaiting for more changes...\n\n", len(flsBackUp))

			} else {
				if deletedString != "" {
					return fmt.Errorf("For some reason it gives deleted values at the start, which should not happen")
				}
				fmt.Printf("[%v✓%v] Files and directories (total in the project: %v) have been added to the watchlist:\n%v\n\n", BrightGreen, Reset, len(flsBackUp), changedString)
			}
		}
		time.Sleep(2000 * time.Millisecond)
	}
}
