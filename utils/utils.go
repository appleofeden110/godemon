package utils

import "log"

// check() takes in error and checks if it's nil or not.
// It takes the boilerplate and basically transforms into a very easy function.
// If needed, error can always be specified, no need to call it on every error in the program
func Check(e error) bool {
	if e != nil {
		log.Fatalf("Some error occured: %v\n", e)
		return true
	}
	return false
}

//
//func Checkf(msg string, e error) bool {
//	if e != nil {
//		log.Fatalf("%v: %v\n", msg, e)
//		return true
//	}
//	return false
//}
