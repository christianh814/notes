package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Since this file doesn't exist...it'll return an error
	_, err := os.Open("no-file.txt")
	if err != nil {
		/* This just prints the error to stdout*/
		//fmt.Println("error happend", err)

		/* This also just prints the error to stdout but with date/time stamp
		   You can also use log.SetOutput(f) to write the output to a file. Note
		   "SetOutput" just sets the output of where log.Println() should write
		   out to. 
		*/
		//log.Println(error happend, err)

		/* This calls os.Exit(1) if it finds and error. This also prints to stdout OR 
		   what is set on log.SetOutput(f). Note that "Fatal" exits immediatly and no
		   defer fuctions are run
		*/
		//log.Fatalln(err)

		/* Panic is like running Println() with panic() afterwards. Panic stops
		   go function and all go subroutines run the deferred actions. Panic
		   can be "recovered"
		*/

		//log.Panicln(err)
	}
}
