package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	dirMap = make(map[string][]string)
	mu     sync.Mutex
	wg     sync.WaitGroup
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
			wg.Add(1)
			go walk(homeDir + entry.Name())
		}
	}

	wg.Wait()
}

func walk(path string) {
	defer wg.Done()

	dirContent, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		return
	}

	for _, entry := range dirContent {
		if entry.IsDir() {
			fmt.Printf("Directory: %s\n", entry.Name())
			wg.Add(1)
			go walk(path + "/" + entry.Name())
		} else {
			fmt.Printf("File: %s\n", entry.Name())
			mu.Lock()
			dirMap[entry.Name()] = append(
				dirMap[entry.Name()],
				filepath.Join(path, entry.Name()),
			)
			mu.Unlock()
		}
	}
}
