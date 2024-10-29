package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Config struct {
	CachePath     string `json:"cachePath"`
	EncryptionKey string `json:"encryptionKey"`
}

var (
	dirMap = make(map[string][]string)
	mu     sync.Mutex
	wg     sync.WaitGroup
)

func main() {
	// Load configuration
	config, err := loadConfig("config.json")
	if err != nil {
		panic(err)
	}

	homeDir := "/"

	rootDir, err := os.ReadDir(homeDir)
	if err != nil {
		panic(err)
	}

	start := time.Now()

	if _, err := os.Stat(config.CachePath); err == nil {
		file, err := os.Open(config.CachePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		ciphertext, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}

		key := []byte(config.EncryptionKey)
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
				wg.Add(1)
				go walk(homeDir + entry.Name())
			}
		}

		wg.Wait()
	}

	elapsed := time.Since(start)
	fmt.Printf("Time taken Stage 1: %s\n", elapsed)

	start = time.Now()
	searchFileName := ".bashrc"

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

	elapsed = time.Since(start)
	fmt.Printf("Time taken Stage 2: %s\n", elapsed)

	start = time.Now()
	file, err := os.Create(config.CachePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	plaintext := new(bytes.Buffer)
	encoder := gob.NewEncoder(plaintext)
	if err := encoder.Encode(dirMap); err != nil {
		panic(err)
	}

	key := []byte(config.EncryptionKey)
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

	elapsed = time.Since(start)
	fmt.Printf("Time taken Stage 3: %s\n", elapsed)
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

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
