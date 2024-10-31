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

func walk(path string) {
	defer wg.Done()

	dirContent, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		return
	}

	for _, entry := range dirContent {
		if entry.IsDir() {
			wg.Add(1)
			go walk(path + "/" + entry.Name())
		} else {
			mu.Lock()
			dirMap[entry.Name()] = append(
				dirMap[entry.Name()],
				filepath.Join(path, entry.Name()),
			)
			mu.Unlock()
		}
	}
}
