package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Message string   `json:"message"`
	Result  []string `json:"result"`
}

func main() {

	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config file")
		return
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Invalid Usage")
		return
	} else {
		switch args[0] {
		case "kill":
			kill(config)
		case "search":
			search(config, args[1])
		case "recache":
			// Recache the cache
		case "status":
			// Get the status of the server
		}
	}
}

func kill(config *Config) {
	url := fmt.Sprintf("http://%s:%d/shutdown", config.Ip, config.Port)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending shutdown request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to shutdown server, status code:", resp.StatusCode)
		return
	}

	fmt.Println("Server shutdown successfully")
}

func search(config *Config, query string) {
	url := fmt.Sprintf("http://%s:%d/search", config.Ip, config.Port)
	payload := fmt.Sprintf(`{"file_name": "%s"}`, query)
	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		fmt.Println("Error sending search request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to search for file, status code:", resp.StatusCode)
		return
	}

	fmt.Println("Search results:")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	fmt.Println("Message:", response.Message)
	fmt.Println("Result Count: ", len(response.Result))
	fmt.Println("Results:")
	for _, result := range response.Result {
		fmt.Println(result)
	}
}
