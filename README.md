# Redis Clone in Go

A lightweight, concurrent, in-memory key-value store built from scratch in Go. This project implements a subset of standard Redis commands and includes a custom parser for the RESP (REdis Serialization Protocol), allowing it to work natively with the standard `redis-cli`.

## Features

- **Standard Commands**: Supports `PING`, `SET`, `GET`, `DEL`, `TTL`, and `EXPIRE`.
- **Key Expiration**: Background cleanup process using a Min-Heap priority queue for efficient TTL management.
- **Eviction Policy**: Implements a Volatile-TTL eviction strategy that removes the soonest-to-expire keys when the database reaches its capacity limit.
- **RESP Protocol**: Custom-built RESP encoder/decoder that handles Simple Strings, Errors, Integers, Bulk Strings, and Arrays.
- **Concurrency**: Thread-safe operations and handling of multiple concurrent client connections using Goroutines.

## Getting Started

### Prerequisites
- [Go](https://golang.org/doc/install) 1.20 or higher.

### Running the Server

Start the server locally (it will listen on port `6379`):

```bash
go run main.go
```

### Connecting to the Server

You can connect to the server using the official `redis-cli` tool:

```bash
redis-cli -p 6379
```

Or you can use `telnet` or `netcat`:

```bash
telnet localhost 6379
```

## Example Usage

```bash
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> SET mykey "Hello World"
OK
127.0.0.1:6379> GET mykey
"Hello World"
127.0.0.1:6379> SET temp "I will expire soon" EX 10
OK
127.0.0.1:6379> TTL temp
(integer) 8
127.0.0.1:6379> DEL mykey
(integer) 1
```

## Architecture

- **`main.go`**: Entry point that initializes the database, the cleanup ticker, and the TCP server.
- **`server/`**: Manages incoming TCP connections and command routing.
- **`internals/db/`**: The core thread-safe map storage, including the Min-Heap based TTL and eviction engine.
- **`internals/protocols/`**: RESP byte-stream parser and encoder.
- **`internals/commands/`**: Handlers for individual Redis commands.
