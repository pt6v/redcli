package redis

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host      string
	Port      int
	Password  string
	Database  int
	Heartbeat int
	Writable  bool
}

type Client struct {
	client    *redis.Client
	config    Config
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	closeChan chan struct{}
}

func NewClient(config Config) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Database,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	c := &Client{
		client:    client,
		config:    config,
		ctx:       ctx,
		cancel:    cancel,
		closeChan: make(chan struct{}),
	}

	// Start heartbeat
	c.wg.Add(1)
	go c.heartbeat()

	return c, nil
}

func (c *Client) heartbeat() {
	defer c.wg.Done()

	ticker := time.NewTicker(time.Duration(c.config.Heartbeat) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeChan:
			return
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.client.Ping(c.ctx).Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Heartbeat failed: %v\n", err)
			}
		}
	}
}

func (c *Client) Execute(cmd string, args ...interface{}) (interface{}, error) {
	// Convert command to Redis command
	cmdArgs := make([]interface{}, 0, len(args)+1)
	cmdArgs = append(cmdArgs, cmd)
	cmdArgs = append(cmdArgs, args...)

	result := c.client.Do(c.ctx, cmdArgs...)
	if result.Err() != nil {
		return nil, result.Err()
	}

	// Parse the result based on command type
	return c.parseResult(cmd, result)
}

func (c *Client) parseResult(cmd string, result *redis.Cmd) (interface{}, error) {
	switch strings.ToUpper(cmd) {
	case "GET":
		return result.Text()
	case "SET":
		return result.Result()
	case "HGET", "HGETALL":
		return c.parseHashResult(cmd, result)
	case "LRANGE", "LLEN":
		return result.StringSlice()
	case "ZRANGE", "ZREVRANGE":
		return c.parseSortedSetResult(result)
	case "SMEMBERS", "SISMEMBER":
		return result.StringSlice()
	case "EXISTS":
		return result.Int()
	case "TTL":
		return result.Int()
	case "TYPE":
		return result.Text()
	case "KEYS":
		return result.StringSlice()
	case "DBSIZE":
		return result.Int()
	case "INFO":
		return result.Text()
	case "PING":
		return result.Text()
	default:
		// Try to return as string for most commands
		return result.Text()
	}
}

func (c *Client) parseHashResult(cmd string, result *redis.Cmd) (interface{}, error) {
	if strings.ToUpper(cmd) == "HGET" {
		return result.Text()
	}

	// HGETALL - get as slice of strings
	vals, err := result.StringSlice()
	if err != nil {
		return nil, err
	}

	// Convert slice to field pairs
	var fields []HashField
	for i := 0; i < len(vals); i += 2 {
		if i+1 < len(vals) {
			fields = append(fields, HashField{
				Key:   vals[i],
				Value: vals[i+1],
			})
		}
	}

	return HashResult{Fields: fields}, nil
}

func (c *Client) parseSortedSetResult(result *redis.Cmd) (interface{}, error) {
	// Try to parse as []interface{}
	vals, err := result.Slice()
	if err != nil {
		return nil, err
	}

	var members []SortedSetMember
	for i := 0; i < len(vals); i += 2 {
		if i+1 < len(vals) {
			members = append(members, SortedSetMember{
				Score:  vals[i],
				Member: vals[i+1],
			})
		}
	}

	return SortedSetResult{Members: members}, nil
}

func (c *Client) Close() error {
	close(c.closeChan)
	c.cancel()
	c.wg.Wait()
	return c.client.Close()
}
