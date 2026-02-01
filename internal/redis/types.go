package redis

// HashField represents a field-value pair in a Redis hash
type HashField struct {
	Key   string
	Value string
}

// HashResult represents the result of an HGETALL command
type HashResult struct {
	Fields []HashField
}

// SortedSetMember represents a member-score pair in a Redis sorted set
type SortedSetMember struct {
	Score  interface{}
	Member interface{}
}

// SortedSetResult represents the result of a ZRANGE command
type SortedSetResult struct {
	Members []SortedSetMember
}

// StringArray represents an array of strings (for lists, sets, etc.)
type StringArray struct {
	Values []string
}

// StatusResult represents a simple status response (OK, etc.)
type StatusResult struct {
	Status string
}

// IntegerResult represents an integer response
type IntegerResult struct {
	Value int64
}
