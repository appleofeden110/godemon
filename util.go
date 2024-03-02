package main

import "log"

// check() takes in error and checks if it's nil or not.
// It takes the boilerplate and basically transforms into a very easy function.
// If needed, error can always be specified, no need to call it on every error in the program
func check(e error) {
	if e != nil {
		log.Fatalf("Some error occured: %v\n", e)
	}
}
