# RedCli - Enhanced Redis Command Line Tool

RedCli 是一个用 Go 语言开发的增强型 Redis 命令行交互工具，旨在提供比官方 redis-cli 更好的使用体验。

## 功能特性

### ✨ 已实现功能

1. **中文显示优化** - 正确显示中文字符，解决 redis-cli 中文编码问题
2. **JSON 格式化** - 自动识别和美化显示 JSON 格式的数据
3. **彩色输出** - 为不同的数据类型和结构添加颜色，提高可读性
4. **只读模式** - 安全模式，禁止执行写操作，防止误操作
5. **心跳机制** - 可配置的心跳检测，保持连接活跃（默认 30 秒）
6. **连接管理** - 支持密码认证、数据库选择等完整连接参数

## 安装和编译

### 前置要求

- Go 1.16 或更高版本
- Redis 服务器（用于连接和测试）

### 编译

```bash
# 克隆或进入项目目录
cd redcli

# 下载依赖
go mod tidy

# 编译
go build -o redcli

# （可选）安装到系统路径
sudo mv redcli /usr/local/bin/
```

## 使用方法

### 基本用法

```bash
# 连接到本地 Redis（默认 127.0.0.1:6379）
./redcli

# 连接到指定主机和端口
./redcli -host 192.168.1.100 -port 6379

# 使用密码认证
./redcli -password yourpassword

# 指定数据库
./redcli -database 1

# 启用写操作（默认只读）
./redcli -writable

# 启用 JSON 美化显示
./redcli -pretty

# 设置心跳间隔（秒）
./redcli -heartbeat 60

# 禁用彩色输出
./redcli -no-color
```

### 完整参数列表

```
-host string
        Redis 服务器地址 (默认 "127.0.0.1")

-port int
        Redis 服务器端口 (默认 6379)

-password string
        Redis 密码

-database int
        Redis 数据库编号 (默认 0)

-writable
        启用写命令（默认：只读模式）

-pretty
        美化显示 JSON 值

-heartbeat int
        心跳间隔（秒） (默认 30)

-no-color
        禁用彩色输出
```

## 使用示例

### 1. 基本操作

```bash
$ ./redcli
Connected to Redis at 127.0.0.1:6379
Running in READ-ONLY mode
Type 'exit' or 'quit' to exit
redcli> PING
PONG
redcli> SET name John
Error: Command 'SET' is not allowed in read-only mode. Use --writable flag to enable write commands.
```

### 2. 中文显示

```bash
redcli> SET username 张三
OK
redcli> GET username
张三
```

### 3. JSON 美化

```bash
redcli> SET user:1 '{"name":"张三","age":25,"city":"北京"}'
OK
redcli> GET user:1 --pretty
{
  "name": "张三",
  "age": 25,
  "city": "北京"
}
```

### 4. Hash 操作（带颜色）

```bash
redcli> HSET user:1000 name "李四" age 30 email "lisi@example.com"
redcli> HGETALL user:1000
field: name
value: 李四

field: age
value: 30

field: email
value: lisi@example.com
```

### 5. Sorted Set 操作

```bash
redcli> ZADD leaderboard 100 "player1" 200 "player2" 150 "player3"
redcli> ZRANGE leaderboard 0 -1 WITHSCORES
score: 100
member: player1

score: 150
member: player3

score: 200
member: player2
```

### 6. 写模式

```bash
# 只读模式（默认）
$ ./redcli
redcli> SET key value
Error: Command 'SET' is not allowed in read-only mode.

# 启用写模式
$ ./redcli -writable
redcli> SET key value
OK
```

## 支持的 Redis 命令

RedCli 支持所有标准 Redis 命令，并针对以下命令进行了特殊优化：

- **String**: GET, SET, MGET, MSET, INCR, DECR 等
- **Hash**: HGET, HSET, HMGET, HMSET, HGETALL, HDEL 等
- **List**: LRANGE, LLEN, LPUSH, RPUSH, LPOP, RPOP 等
- **Set**: SMEMBERS, SADD, SREM, SISMEMBER 等
- **Sorted Set**: ZRANGE, ZREVRANGE, ZADD, ZREM 等
- **Key**: KEYS, EXISTS, TYPE, TTL, EXPIRE, DEL 等
- **Server**: PING, INFO, DBSIZE 等

## 只读模式

在只读模式（默认）下，以下命令将被禁止：

- 所有写操作命令：SET, DEL, HSET, LPUSH, ZADD 等
- 所有删除操作命令：DEL, UNLINK, HDEL 等
- 所有过期时间设置命令：EXPIRE, EXPIREAT 等

要启用写操作，请使用 `-writable` 参数。

## 项目结构

```
redcli/
├── main.go                 # 主程序入口
├── go.mod                  # Go 模块文件
├── go.sum                  # 依赖锁定文件
├── internal/
│   ├── redis/             # Redis 客户端实现
│   │   ├── client.go      # 客户端核心逻辑
│   │   └── types.go       # 数据类型定义
│   ├── command/           # 命令解析和处理
│   │   └── parser.go      # 命令解析器
│   └── display/           # 显示和格式化
│       └── display.go     # 显示逻辑
├── README.md              # 项目说明文档
├── 需求列表.md            # 需求文档
└── 优化建议.md            # 未来优化建议
```

## 与 redis-cli 的对比

| 特性 | redis-cli | redcli |
|------|-----------|--------|
| 中文显示 | 需要设置 --raw | 原生支持 |
| JSON 格式化 | 不支持 | 支持 |
| 彩色输出 | 不支持 | 支持 |
| 只读模式 | 不支持 | 支持 |
| 心跳检测 | 不支持 | 支持（可配置） |
| 交互体验 | 基础 | 增强型 |

## 未来优化方向

详见 [优化建议.md](tmp/优化建议.md) 文档，包含以下方向：

1. 命令历史和自动补全
2. 配置文件支持
3. 批量操作和脚本支持
4. 数据导出和导入功能
5. 监控和性能分析
6. 多数据库支持增强
7. 事务支持
8. SSH 隧道支持
9. 插件系统
10. TLS/SSL 支持
11. 集群模式支持
12. 键空间分析
13. 安全增强

## 常见问题

### Q: 如何连接到需要密码的 Redis？

A: 使用 `-password` 参数：

```bash
./redcli -password yourpassword
```

### Q: 如何在脚本中使用 redcli？

A: 可以通过管道传递命令：

```bash
echo -e "PING\nGET mykey" | ./redcli
```

### Q: 中文显示还是乱码怎么办？

A: 确保您的终端支持 UTF-8 编码。大多数现代终端默认支持 UTF-8。

### Q: 为什么默认是只读模式？

A: 只读模式可以防止误操作导致的数据丢失或损坏。如果您需要执行写操作，请明确使用 `-writable` 参数。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 作者

Claude Code

## 致谢

本项目使用了以下优秀的 Go 库：

- [go-redis](https://github.com/redis/go-redis) - Redis 客户端库
