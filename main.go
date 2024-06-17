package main

import (
	"fmt"
	"github.com/appleofeden110/godemon/shell"
	"github.com/appleofeden110/godemon/tree"
	"log"
	"time"
)

func mapCompare(map1, map2 map[tree.KeyFile]time.Time) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, val := range map1 {
		if val2, exists := map2[key]; !exists || !val.Equal(val2) {
			return false
		}
	}
	return true
}

func GodemonInit() error {
	flsBackUp := make(map[tree.KeyFile]time.Time)
	fls := make(map[tree.KeyFile]time.Time)

	for {
		_, err := tree.BLR(".", fls)
		check(err, "tree problem")
		if !mapCompare(flsBackUp, fls) {
			//for now, no shell()

			// Create a new map and deep copy fls into it
			newBackup := make(map[tree.KeyFile]time.Time)
			for k, v := range fls {
				newBackup[k] = v
			}
			flsBackUp = newBackup // Now flsBackUp is a deep copy of fls
			log.Println("size: ", len(flsBackUp))
			log.Println(flsBackUp)
			log.Println("Files have changed, action taken.")
		} else {
			log.Println("No changes detected.")
		}

		fmt.Println("Waiting for changes...")
		time.Sleep(400 * time.Millisecond)
	}
}

func check(err error, msg ...string) {
	if err != nil {
		log.Fatalf("GODEMON ERR: %s", msg)
	}
}

func main() {
	err := GodemonInit()
	if err != nil {
		panic(err)
	}
}

func ShellHAHA() (*shell.File, error) {
	return shell.CreateFile()
}
