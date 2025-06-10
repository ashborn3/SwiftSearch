package main

import (
	"context"
	"fmt"
	"log"
	"swift_search/internal/cache"
	"swift_search/internal/config"
	"swift_search/internal/router"
)

type Server struct {
	Config *config.Config
}

func main() {
	fmt.Printf("Server Started!\n")

	server, err := InitServer()
	if err != nil {
		panic(err)
	}

	err = cache.DeserializeCache(server.Config)
	if err != nil {
		log.Printf("Error during deserialization: %s", err.Error())
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go cache.SyncCacheToDisk(ctx, server.Config)

	router.Server(ctx, cancel, server.Config)
}

func InitServer() (*Server, error) {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return nil, err
	}

	return &Server{Config: config}, nil
}
