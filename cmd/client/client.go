package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"swift_search/internal/config"
)

type Client struct {
	Config *config.Config
}

type Response struct {
	Message string   `json:"message"`
	Result  []string `json:"result"`
}

func main() {
	client, err := InitClient()
	if err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Invalid Usage")
		return
	} else {
		switch args[0] {
		case "kill":
			client.kill()
		case "search":
			client.search(args[1])
		case "recache":
			client.recache()
		case "status":
			client.status()
		default:
			fmt.Println("Invalid Command")
		}
	}
}

func InitClient() (*Client, error) {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config file")
		return nil, err
	}

	return &Client{Config: config}, nil
}

func (cl *Client) kill() {
	url := fmt.Sprintf("http://%s:%d/shutdown", cl.Config.Ip, cl.Config.Port)
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

func (cl *Client) search(query string) {
	url := fmt.Sprintf("http://%s:%d/search", cl.Config.Ip, cl.Config.Port)
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

func (cl *Client) recache() {
	url := fmt.Sprintf("http://%s:%d/recache", cl.Config.Ip, cl.Config.Port)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending recache request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to recache files, status code:", resp.StatusCode)
		return
	}

	fmt.Println("Files recached successfully")
}

func (cl *Client) status() {
	url := fmt.Sprintf("http://%s:%d/status", cl.Config.Ip, cl.Config.Port)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending status request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get server status, status code:", resp.StatusCode)
		return
	}

	fmt.Println("Server status: OK")
}
