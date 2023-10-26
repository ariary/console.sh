package check

import (
	"fmt"
	"os"
)

//Check: Check an error and display appropriated message (don't exit or panic). Print the error if msg is empty.
func Check(e error, msg string) {
	if e != nil {
		if msg != "" {
			fmt.Println(msg)
		}
		fmt.Println(e)
	}
}

//CheckAndExit: Check an error and display appropriated message, and then exit
func CheckAndExit(e error, msg string) {
	if e != nil {
		if msg != "" {
			fmt.Println(msg)
		}
		fmt.Println(e)
		os.Exit(1)
	}
}
