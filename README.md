# Simple Redis Client in Go

Redis-go is a client for the [redis](https://github.com/antirez/redis) Key-Value Storage system.

## Example

```go
package main

import "github.com/zonyitoo/redis-go"

func main() {
    client := redis.NewClient("127.0.0.1:6379")
    _, _ := client.Exec("set", "hello", "world")

    ret, _ := client.Exec("get", "hello")

    println(ret.(redis.RespBulkString).data)
}

```
