package redis

import (
    "bufio"
    "errors"
    "net"
)

type Client struct {
    conn   net.Conn
    parser *RespParser
    writer *bufio.Writer
}

func NewClient(ipaddr string) *Client {
    conn, err := net.Dial("tcp", ipaddr)
    if err != nil {
        panic(err)
    }

    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

    return &Client{
        conn:   conn,
        parser: NewRespParser(reader),
        writer: writer,
    }
}

func (c *Client) Exec(cmdname string, args ...string) (RespObject, error) {
    cmds := RespList{
        data: []RespObject{RespBulkString{data: cmdname}},
    }
    for _, cmd := range args {
        cmds.data = append(cmds.data, RespBulkString{data: cmd})
    }
    _, err := c.writer.WriteString(cmds.ToResp())
    if err != nil {
        return nil, err
    }
    c.writer.Flush()

    obj, err := c.parser.next()
    if err != nil {
        return nil, err
    }

    switch t := obj.(type) {
    case RespErrorString:
        return nil, errors.New(t.message)
    default:
        return obj, nil
    }
}
