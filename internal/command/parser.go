package command

import (
	"strings"
)

// List of write commands that should be blocked in read-only mode
var writeCommands = map[string]bool{
	"SET": true,
	"SETNX": true,
	"SETEX": true,
	"PSETEX": true,
	"MSET": true,
	"MSETNX": true,
	"GETSET": true,
	"APPEND": true,
	"SETRANGE": true,
	"INCR": true,
	"INCRBY": true,
	"INCRBYFLOAT": true,
	"DECR": true,
	"DECRBY": true,
	"DEL": true,
	"UNLINK": true,
	"EXPIRE": true,
	"EXPIREAT": true,
	"PEXPIRE": true,
	"PEXPIREAT": true,
	"PERSIST": true,
	"HSET": true,
	"HSETNX": true,
	"HMSET": true,
	"HINCRBY": true,
	"HINCRBYFLOAT": true,
	"HDEL": true,
	"LPUSH": true,
	"RPUSH": true,
	"LPOP": true,
	"RPOP": true,
	"LINSERT": true,
	"LSET": true,
	"LTRIM": true,
	"RPOPLPUSH": true,
	"SADD": true,
	"SREM": true,
	"SPOP": true,
	"SMOVE": true,
	"SINTERSTORE": true,
	"SUNIONSTORE": true,
	"SDIFFSTORE": true,
	"ZADD": true,
	"ZINCRBY": true,
	"ZREM": true,
	"ZREMRANGEBYRANK": true,
	"ZREMRANGEBYSCORE": true,
	"ZREMRANGEBYLEX": true,
	"ZUNIONSTORE": true,
	"ZINTERSTORE": true,
	"PFADD": true,
	"PFMERGE": true,
	"PUBLISH": true,
	"RENAME": true,
	"RENAMENX": true,
	"MIGRATE": true,
	"MOVE": true,
	"RESTORE": true,
	"DUMP": true,
	"FLUSHDB": true,
	"FLUSHALL": true,
	"SORT": true,
	"BITOP": true,
	"GEOADD": true,
}

// Parse parses the input string and returns the command and arguments
func Parse(input string) (string, []interface{}) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}

	cmd := strings.ToUpper(parts[0])
	args := make([]interface{}, 0, len(parts)-1)

	for i := 1; i < len(parts); i++ {
		args = append(args, parts[i])
	}

	return cmd, args
}

// IsWriteCommand checks if a command is a write command
func IsWriteCommand(cmd string) bool {
	return writeCommands[strings.ToUpper(cmd)]
}

// IsReadCommand checks if a command is a read command
func IsReadCommand(cmd string) bool {
	return !IsWriteCommand(cmd)
}
