package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func server(ctx context.Context, cancel context.CancelFunc, config *Config) {
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is running",
		})
	})

	r.GET("/recache", func(c *gin.Context) {
		log.Printf("Recaching...")
		err := deserializeCache(config)
		if err != nil {
			log.Printf("Recaching Failed: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Error recaching: %s", err),
			})
		} else {
			log.Printf("Recaching Success")
			c.JSON(http.StatusOK, gin.H{
				"message": "Cache Recached",
			})
		}
	})

	r.GET("/shutdown", func(c *gin.Context) {
		log.Println("Shutting down server...")
		cancel()
		c.JSON(http.StatusOK, gin.H{
			"message": "Shutting down server",
		})
	})

	r.POST("/search", func(c *gin.Context) {
		var request struct {
			FileName string `json:"file_name"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
			return
		}

		log.Println("Received file name:", request.FileName)

		resultArr, exists := dirMap[request.FileName]
		if exists {
			log.Printf("File Found\n")
			c.JSON(http.StatusOK, gin.H{
				"message": "Found",
				"result":  resultArr,
			})
		} else {
			log.Printf("File Not Found\n")
			c.JSON(http.StatusOK, gin.H{
				"message": "Not Found",
				"result":  []string{},
			})
		}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Ip, config.Port),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
