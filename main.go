package main

import (
	"github.com/appleofeden110/godemon/tree_check"
)

func main() {
	err := tree_check.TreeCheck()
	if err != nil {
		panic(err)
	}
}
