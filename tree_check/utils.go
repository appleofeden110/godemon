package tree_check

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"time"
)

var (
	ErrNoDirectoriesToCheck = errors.New("All files are ignored, check ignored.json")
	ErrIgnoreDirs           = errors.New("ignoreDirs.json should be created in .godemon (It will be automatized later)")
)

const (
	BrightGreen  = "\033[92m"
	BrightRed    = "\033[91m"
	BrightYellow = "\033[93m"
	Reset        = "\033[0m"
)

func check(err error, msg ...string) {
	if err != nil {
		log.Fatalf("GODEMON ERR: %s: %v", msg, err)
	}
}

func ignoreDirs(ignoreDirs map[string]bool) error {
	jsonF, err := os.Open("ignored.json")
	if err != nil {
		return ErrIgnoreDirs
	}
	defer jsonF.Close()

	b, err := io.ReadAll(jsonF)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	err = json.Unmarshal(b, &ignoreDirs)
	if err != nil {
		return fmt.Errorf("There is an error unmarshaling data: %v\n", err)
	}
	return nil
}

func mapCompare(map1, map2 map[KeyFile]time.Time) ([]KeyFile, []KeyFile, bool) {
	var changedFiles []KeyFile
	var deletedFiles []KeyFile

	for key, val2 := range map2 {
		if !CheckFile(key) {
			DeleteIndex(map2, key)
			deletedFiles = append(deletedFiles, key)
			continue
		}
		if val1, exists := map1[key]; !exists || !val1.Equal(val2) {
			changedFiles = append(changedFiles, key)
		}
	}

	return changedFiles, deletedFiles, len(changedFiles) == 0 && len(deletedFiles) == 0
}

func printFileArray(changedFiles, deletedFiles []KeyFile) (changed string, deleted string) {
	var fileColumn string
	var deletedColumn string

	for _, name := range changedFiles {
		fileColumn += "\n\t- " + name.Name + " (" + name.Path + ")"
	}
	for _, name := range deletedFiles {
		deletedColumn += "\n\t- " + name.Name + " (" + name.Path + ")"
	}
	return fileColumn, deletedColumn
}

func CheckFile(keyFile KeyFile) bool {
	_, err := os.Stat(keyFile.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func DeleteIndex(keyFiles map[KeyFile]time.Time, fileName KeyFile) bool {

	maps.DeleteFunc[map[KeyFile]time.Time, KeyFile, time.Time](keyFiles, func(k KeyFile, v time.Time) bool {
		if k == fileName {
			return true
		}
		return false
	})
	return false
}
