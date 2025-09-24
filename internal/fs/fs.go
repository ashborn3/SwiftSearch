package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileMeta struct {
	Path    string
	ModTime int64
}

var (
	DirMap      = make(map[string][]FileMeta)
	Mu          sync.Mutex
	workerCount = 32 // Number of concurrent workers (increased)
	excludeDirs = map[string]struct{}{
		"mnt": {}, "Windows": {}, ".git": {}, "node_modules": {}, "__pycache__": {}, "venv": {}, ".cache": {},
	}
)

type job struct {
	path string
}

func Walk(root string) {
	var wg sync.WaitGroup
	walkDir(root, &wg)
	wg.Wait()
}

func walkDir(path string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		done := make(chan struct{})
		var dirContent []os.DirEntry
		var err error
		go func() {
			dirContent, err = os.ReadDir(path)
			close(done)
		}()
		select {
		case <-done:
			if err != nil {
				// fmt.Printf("Error reading directory: %s\n", err)
				return
			}
		case <-time.After(1 * time.Second):
			fmt.Printf("Timeout reading directory: %s\n", path)
			return
		}
		// fmt.Printf("Scanning: %s, entries: %d\n", path, len(dirContent))
		for _, entry := range dirContent {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			// Skip symlinks
			if info.Mode()&os.ModeSymlink != 0 {
				continue
			}
			if entry.IsDir() {
				if _, skip := excludeDirs[entry.Name()]; skip {
					continue
				}
				walkDir(filepath.Join(path, entry.Name()), wg)
			} else {
				meta := FileMeta{
					Path:    filepath.Join(path, entry.Name()),
					ModTime: info.ModTime().Unix(),
				}
				Mu.Lock()
				// Check if file already exists and mod time matches
				exists := false
				for i, fm := range DirMap[entry.Name()] {
					if fm.Path == meta.Path {
						exists = true
						if fm.ModTime != meta.ModTime {
							DirMap[entry.Name()][i] = meta // update mod time
						}
						break
					}
				}
				if !exists {
					DirMap[entry.Name()] = append(DirMap[entry.Name()], meta)
				}
				Mu.Unlock()
			}
		}
	}()
}

func scanDir(path string, jobs chan<- job) {
	done := make(chan struct{})
	var dirContent []os.DirEntry
	var err error
	go func() {
		dirContent, err = os.ReadDir(path)
		close(done)
	}()
	select {
	case <-done:
		if err != nil {
			fmt.Printf("Error reading directory: %s\n", err)
			return
		}
	case <-time.After(1 * time.Second):
		fmt.Printf("Timeout reading directory: %s\n", path)
		return
	}
	fmt.Printf("Scanning: %s, entries: %d\n", path, len(dirContent))
	for _, entry := range dirContent {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		// Skip symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if entry.IsDir() {
			if _, skip := excludeDirs[entry.Name()]; skip {
				continue
			}
			jobs <- job{path: filepath.Join(path, entry.Name())}
		} else {
			meta := FileMeta{
				Path:    filepath.Join(path, entry.Name()),
				ModTime: info.ModTime().Unix(),
			}
			Mu.Lock()
			// Check if file already exists and mod time matches
			exists := false
			for i, fm := range DirMap[entry.Name()] {
				if fm.Path == meta.Path {
					exists = true
					if fm.ModTime != meta.ModTime {
						DirMap[entry.Name()][i] = meta // update mod time
					}
					break
				}
			}
			if !exists {
				DirMap[entry.Name()] = append(DirMap[entry.Name()], meta)
			}
			Mu.Unlock()
		}
	}
}
