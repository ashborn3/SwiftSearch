package main

import (
	"context"
	"fmt"
	"swift_search/internal/config"
)

func main() {
	fmt.Printf("Server Started!\n")

	config, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	deserializeCache(config)

	ctx, cancel := context.WithCancel(context.Background())

	go syncCacheToDisk(ctx, config)

	server(ctx, cancel, config)
}
