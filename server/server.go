package main

import (
	"context"
	"fmt"
)

func main() {
	fmt.Printf("Server Started!\n")

	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	fmt.Println(config)

	deserializeCache(config)

	ctx, cancel := context.WithCancel(context.Background())

	go syncCacheToDisk(ctx, config)

	server(ctx, cancel, config)
}
