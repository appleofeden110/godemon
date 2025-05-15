package godemon

import (
	"fmt"
	"log"
	"time"
)

func check(err error, msg ...string) {
	if err != nil {
		log.Fatalf("GODEMON ERR: %s: %v", msg, err)
	}
}

func mapCompare(map1, map2 map[KeyFile]time.Time) bool {
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
	flsBackUp := make(map[KeyFile]time.Time)
	fls := make(map[KeyFile]time.Time)

	for {
		_, err := BLR(".", fls)
		if err != nil {
			return fmt.Errorf("Tree problem: %v\n", err)
		}
		if !mapCompare(flsBackUp, fls) {
			newBackup := make(map[KeyFile]time.Time)
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
