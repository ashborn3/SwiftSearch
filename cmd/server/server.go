package main

import (
	"context"
	"fmt"
	"swift_search/internal/config"
)

type Server struct {
	Config *config.Config
}

func main() {
	fmt.Printf("Server Started!\n")

	deserializeCache(config)

	ctx, cancel := context.WithCancel(context.Background())

	go syncCacheToDisk(ctx, config)

	server(ctx, cancel, config)
}

func InitServer() (*Server, error) {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return nil, err
	}

	return &Server{Config: config}, nil
}
