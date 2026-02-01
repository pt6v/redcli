package display

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"redcli/internal/redis"
)

type Config struct {
	Pretty  bool
	NoColor bool
}

// ANSI color codes
const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	Bold      = "\033[1m"
)

var (
	colorKey     = Cyan
	colorString  = Green
	colorNumber  = Yellow
	colorBool    = Blue
	colorNil     = Magenta
	colorError   = Red
	colorField   = Yellow
	colorMember  = Green
	colorScore   = Yellow
)

// Result displays the result with proper formatting and colors
func Result(result interface{}, config Config) {
	if result == nil {
		fmt.Println("(nil)")
		return
	}

	switch r := result.(type) {
	case string:
		displayString(r, config)
	case int64:
		displayInteger(r, config)
	case []string:
		displayStringArray(r, config)
	case redis.HashResult:
		displayHash(r, config)
	case redis.SortedSetResult:
		displaySortedSet(r, config)
	default:
		// Try to format as JSON for unknown types
		displayAsJSON(result, config)
	}
}

func displayString(s string, config Config) {
	// Check if it's JSON
	if config.Pretty && isJSON(s) {
		displayAsJSONPretty(s, config)
		return
	}

	// Display with proper encoding (handles Chinese characters correctly)
	if !config.NoColor {
		fmt.Printf("%s%s%s\n", colorString, s, Reset)
	} else {
		fmt.Println(s)
	}
}

func displayInteger(n int64, config Config) {
	if !config.NoColor {
		fmt.Printf("%s%d%s\n", colorNumber, n, Reset)
	} else {
		fmt.Printf("%d\n", n)
	}
}

func displayStringArray(arr []string, config Config) {
	if len(arr) == 0 {
		fmt.Println("(empty array)")
		return
	}

	for i, item := range arr {
		if !config.NoColor {
			fmt.Printf("%s%d)%s %s%s%s\n", colorNumber, i+1, Reset, colorString, item, Reset)
		} else {
			fmt.Printf("%d) %s\n", i+1, item)
		}
	}
}

func displayHash(hash redis.HashResult, config Config) {
	if len(hash.Fields) == 0 {
		fmt.Println("(empty hash)")
		return
	}

	for _, field := range hash.Fields {
		if !config.NoColor {
			fmt.Printf("%sfield:%s %s%s%s\n", colorField, Reset, colorKey, field.Key, Reset)
			fmt.Printf("%svalue:%s ", colorField, Reset)

			// Check if value is JSON
			if config.Pretty && isJSON(field.Value) {
				displayAsJSONPretty(field.Value, config)
			} else {
				fmt.Printf("%s%s%s\n", colorString, field.Value, Reset)
			}
		} else {
			fmt.Printf("field: %s\n", field.Key)
			if config.Pretty && isJSON(field.Value) {
				displayAsJSONPretty(field.Value, config)
			} else {
				fmt.Printf("value: %s\n", field.Value)
			}
		}
		fmt.Println()
	}
}

func displaySortedSet(sortedSet redis.SortedSetResult, config Config) {
	if len(sortedSet.Members) == 0 {
		fmt.Println("(empty sorted set)")
		return
	}

	for _, member := range sortedSet.Members {
		if !config.NoColor {
			fmt.Printf("%sscore:%s %s%v%s\n", colorScore, Reset, colorNumber, member.Score, Reset)
			fmt.Printf("%smember:%s %s%v%s\n", colorMember, Reset, colorString, member.Member, Reset)
		} else {
			fmt.Printf("score: %v\n", member.Score)
			fmt.Printf("member: %v\n", member.Member)
		}
		fmt.Println()
	}
}

func displayAsJSONPretty(jsonStr string, config Config) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		// Not valid JSON, display as string
		if !config.NoColor {
			fmt.Printf("%s%s%s\n", colorString, jsonStr, Reset)
		} else {
			fmt.Println(jsonStr)
		}
		return
	}

	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(jsonStr)
		return
	}

	if !config.NoColor {
		fmt.Printf("%s%s%s\n", colorString, string(formatted), Reset)
	} else {
		fmt.Println(string(formatted))
	}
}

func displayAsJSON(data interface{}, config Config) {
	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", data)
		return
	}

	if !config.NoColor {
		fmt.Printf("%s%s%s\n", colorString, string(formatted), Reset)
	} else {
		fmt.Println(string(formatted))
	}
}

func isJSON(s string) bool {
	s = strings.TrimSpace(s)
	return (strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")) ||
		(strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]"))
}

// Error displays an error message
func Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s%s%s", colorError, fmt.Sprintf(format, args...), Reset)
}

// Success displays a success message
func Success(format string, args ...interface{}) {
	fmt.Printf("%s%s%s", colorString, fmt.Sprintf(format, args...), Reset)
}
