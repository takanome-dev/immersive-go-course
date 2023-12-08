package cmd

import (
	"log"
	"os"
)

func Execute() {
  if len(os.Args) < 2 {
    log.Fatal("Missing filename argument")
  }

  filename := os.Args[1]

  info, err := os.Stat(filename)

  if os.IsNotExist(err) {
    log.Fatal("File does not exist")
  }

  if info.IsDir() {
    log.Fatal("Cannot use directories, please provide a filename")
  }

  fileContent, err := os.ReadFile(filename)
  if err != nil {
    log.Fatal(err)
  }

  os.Stdout.Write(fileContent)
}


