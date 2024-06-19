package main

import (
	"github.com/appleofeden110/godemon/init"
	"log"
)

func check(err error, msg ...string) {
	if err != nil {
		log.Fatalf("GODEMON ERR: %s: %v", msg, err)
	}
}

func main() {
	err := init.GodemonInit()
	if err != nil {
		panic(err)
	}
}
