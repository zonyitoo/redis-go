package redis

import (
    "testing"
)

var client *Client = NewClient("127.0.0.1:6379")

func TestBasic(t *testing.T) {
    ret, err := client.Exec("set", "hello", "world")
    if err != nil {
        t.Fatal(err)
    }
    t.Log(ret)

    _, err = client.Exec("set", "test")
    if err == nil {
        t.Fatal("Should failed because of argument number error")
    }

    ret, err = client.Exec("get", "hello")
    if err != nil {
        t.Fatal(err)
    }
    if ret.(RespBulkString).data != "world" {
        t.Fail()
    }
}
