package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) > 1 {
		fmt.Println("Ahh, unfortunately I can currently only compile one file at a time..")
		fmt.Println("If you would like to compile more than one file, please add import statements to your main file")
		fmt.Println("like so: import 'test'")
		fmt.Println("Thanks!")
		fmt.Println(" ")
		os.Exit(1)
	} else {
		Parser(ReadInFile(args[0]))
	}
}
