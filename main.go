package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"redcli/internal/command"
	"redcli/internal/display"
	"redcli/internal/redis"
)

var (
	host       = flag.String("host", "127.0.0.1", "Redis server host")
	port       = flag.Int("port", 6379, "Redis server port")
	password   = flag.String("password", "", "Redis server password")
	database   = flag.Int("database", 0, "Redis database number")
	writable   = flag.Bool("writable", false, "Enable write commands (default: read-only mode)")
	pretty     = flag.Bool("pretty", false, "Pretty print JSON values")
	heartbeat  = flag.Int("heartbeat", 30, "Heartbeat interval in seconds")
	noColor    = flag.Bool("no-color", false, "Disable colored output")
)

func main() {
	flag.Parse()

	// Create Redis client configuration
	config := redis.Config{
		Host:       *host,
		Port:       *port,
		Password:   *password,
		Database:   *database,
		Heartbeat:  *heartbeat,
		Writable:   *writable,
	}

	// Connect to Redis
	client, err := redis.NewClient(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	fmt.Printf("Connected to Redis at %s:%d\n", *host, *port)
	if !*writable {
		fmt.Println("Running in READ-ONLY mode")
	}
	fmt.Println("Type 'exit' or 'quit' to exit")

	// Setup display config
	displayConfig := display.Config{
		Pretty:  *pretty,
		NoColor: *noColor,
	}

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		client.Close()
		os.Exit(0)
	}()

	// Read-eval-print loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("redcli> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle exit commands
		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			fmt.Println("Exiting...")
			break
		}

		// Parse command
		cmd, args := command.Parse(input)
		if cmd == "" {
			continue
		}

		// Check if command is allowed in read-only mode
		if !*writable && command.IsWriteCommand(cmd) {
			display.Error("Command '%s' is not allowed in read-only mode. Use --writable flag to enable write commands.\n", cmd)
			continue
		}

		// Execute command
		result, err := client.Execute(cmd, args...)
		if err != nil {
			display.Error("Error: %v\n", err)
			continue
		}

		// Display result
		display.Result(result, displayConfig)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
