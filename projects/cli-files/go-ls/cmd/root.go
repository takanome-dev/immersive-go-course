package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func Execute() {
	var help bool

	flag.BoolVar(&help, "h", false, "help")
	flag.Parse()

	if help {
		fmt.Println("Usage: ls [path]")
		fmt.Println("go-ls is a simple implementation of the ls command in Go")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		currentDir = os.Args[1]
	}

	file, err := os.Stat(currentDir)
	if err != nil {
		log.Fatalf("an error occurred, %v", err)
		os.Exit(1)
	}

	if file.IsDir() {
		printDirContents(currentDir)
	} else {
		fmt.Fprintf(os.Stdout, file.Name())
	}
}

func printDirContents(path string) {
	dir, err := os.Open(path)
	if err != nil {
		log.Fatalf("an error occurred, %v", err)
		os.Exit(1)
	}

	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatalf("an error occurred, %v", err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Fprintf(os.Stdout, file.Name()+"\n")
	}
}
