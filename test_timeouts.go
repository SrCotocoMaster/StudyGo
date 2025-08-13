package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("Testing timeout configurations...")

	// Start the server in background
	fmt.Println("Starting server...")
	cmd := exec.Command("go", "run", "server.go")
	err := cmd.Start()
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer cmd.Process.Kill()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Test API endpoint
	fmt.Println("Testing /cotacao endpoint...")
	resp, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		log.Printf("Error calling API: %v", err)
	} else {
		fmt.Printf("API response status: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	// Run client
	fmt.Println("Running client...")
	clientCmd := exec.Command("go", "run", "client.go")
	output, err := clientCmd.CombinedOutput()
	if err != nil {
		log.Printf("Client error: %v", err)
	}
	fmt.Printf("Client output: %s\n", string(output))

	fmt.Println("Test completed!")
}
