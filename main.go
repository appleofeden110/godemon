package main

import godemon "github.com/appleofeden110/godemon/tree-check"

func main() {
	err := godemon.GodemonInit()
	if err != nil {
		panic(err)
	}
}
