package cache

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"swift_search/internal/config"
	"swift_search/internal/fs"
	"time"
)

func DeserializeCache(config *config.Config) error {
	rootDir, err := os.ReadDir(config.HomePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(config.CachePath); err == nil {
		file, err := os.Open(config.CachePath)
		if err != nil {
			return err
		}
		defer file.Close()

		ciphertext, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		key := []byte(config.EncryptionKey)
		block, err := aes.NewCipher(key)
		if err != nil {
			return err
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return err
		}

		nonceSize := gcm.NonceSize()
		if len(ciphertext) < nonceSize {
			return fmt.Errorf("ciphertext too short")
		}

		nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return err
		}

		buffer := bytes.NewBuffer(plaintext)
		decoder := gob.NewDecoder(buffer)
		if err := decoder.Decode(&fs.DirMap); err != nil {
			return err
		}
	} else {
		for _, entry := range rootDir {
			if entry.IsDir() && entry.Name() != "mnt" && entry.Name() != "Windows" {
				fs.Wg.Add(1)
				go fs.Walk(config.HomePath + entry.Name())
			}
		}

		fs.Wg.Wait()
	}

	return nil
}

func SerializeCache(config *config.Config) error {
	file, err := os.Create(config.CachePath)
	if err != nil {
		return err
	}
	defer file.Close()

	plaintext := new(bytes.Buffer)
	encoder := gob.NewEncoder(plaintext)
	if err := encoder.Encode(fs.DirMap); err != nil {
		return err
	}

	key := []byte(config.EncryptionKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext.Bytes(), nil)
	if _, err := file.Write(ciphertext); err != nil {
		return err
	}

	return nil
}

func SyncCacheToDisk(ctx context.Context, config *config.Config) {
	cacheTicker := time.NewTicker(time.Duration(config.SyncTime) * time.Minute)
	defer cacheTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-cacheTicker.C:
			if err := SerializeCache(config); err != nil {
				fmt.Printf("Error serializing cache: %v\n", err)
			}
		}
	}
}
