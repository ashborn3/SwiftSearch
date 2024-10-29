package main

import (
	"fmt"
	"os"
)

func main() {
	homeDir := "/"

	rootDir, err := os.ReadDir(homeDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range rootDir {
		if entry.IsDir() && entry.Name() != "mnt" {
			fmt.Printf("Directory: %s\n", entry.Name())
			walk(homeDir + entry.Name())
		}
	}
}

func walk(path string) {
	dirContent, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		return
	}

	for _, entry := range dirContent {
		if entry.IsDir() {
			fmt.Printf("Directory: %s\n", entry.Name())
			// walk(path + "/" + entry.Name())
		} else {
			fmt.Printf("File: %s\n", entry.Name())
			// Do something with the file
		}
	}
}
