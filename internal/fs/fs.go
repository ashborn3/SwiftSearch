package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	DirMap = make(map[string][]string)
	Mu     sync.Mutex
	Wg     sync.WaitGroup
)

func Walk(path string) {
	defer Wg.Done()

	dirContent, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		return
	}

	for _, entry := range dirContent {
		if entry.IsDir() {
			Wg.Add(1)
			go Walk(path + "/" + entry.Name())
		} else {
			Mu.Lock()
			DirMap[entry.Name()] = append(
				DirMap[entry.Name()],
				filepath.Join(path, entry.Name()),
			)
			Mu.Unlock()
		}
	}
}
