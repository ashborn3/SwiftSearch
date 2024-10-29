package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
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

	if _, err := os.Stat("cache.gob"); err == nil {
		file, err := os.Open("dirMap.gob")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		ciphertext, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}

		key := []byte("a very very very very secret key") // 32 bytes for AES-256
		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err)
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err)
		}

		nonceSize := gcm.NonceSize()
		if len(ciphertext) < nonceSize {
			panic("ciphertext too short")
		}

		nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			panic(err)
		}

		buffer := bytes.NewBuffer(plaintext)
		decoder := gob.NewDecoder(buffer)
		if err := decoder.Decode(&dirMap); err != nil {
			panic(err)
		}
	} else {
		for _, entry := range rootDir {
			if entry.IsDir() && entry.Name() != "mnt" {
				fmt.Printf("Directory: %s\n", entry.Name())
				wg.Add(1)
				go walk(homeDir + entry.Name())
			}
		}

		wg.Wait()
	}

	searchFileName := "targetFileName"

	mu.Lock()
	paths, found := dirMap[searchFileName]
	mu.Unlock()

	if found {
		fmt.Printf("File %s found at:\n", searchFileName)
		for _, path := range paths {
			fmt.Println(path)
		}
	} else {
		fmt.Printf("File %s not found\n", searchFileName)
	}

	file, err := os.Create("cache.gob")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	plaintext := new(bytes.Buffer)
	encoder := gob.NewEncoder(plaintext)
	if err := encoder.Encode(dirMap); err != nil {
		panic(err)
	}

	key := []byte("a very very very very secret key") // 32 bytes for AES-256
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext.Bytes(), nil)
	if _, err := file.Write(ciphertext); err != nil {
		panic(err)
	}
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
