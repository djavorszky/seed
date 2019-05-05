package main

import (
	"flag"
	"fmt"
	"os"
	"seed"
)

var (
	initialize bool
	name       string
)

func main() {
	flag.BoolVar(
		&initialize, "i", false, "Specify this flag to initialize a project.")
	flag.StringVar(&name, "n", "example2", "Specify the project's name. A new folder will be created "+
		"with this name where the poject will be initialized. If '.' is specified, the project name will be derived "+
		"from the directory name, and the project will be initialized in the same folder.")

	flag.Parse()

	if !initialize {
		fmt.Println("debug: init flag not set, please set it to continue testing.")
		flag.Usage()
		os.Exit(1)
	}

	err := seed.InitProject(name)
	if err != nil {
		fmt.Printf("Failed initializing the project: %v\n", err)
		os.Exit(1)
	}
}
