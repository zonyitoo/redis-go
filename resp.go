package redis

import (
    "bufio"
    "errors"
    "fmt"
    "strconv"
    "strings"
)

type RespParser struct {
    bufr *bufio.Reader
}

func NewRespParser(bufr *bufio.Reader) *RespParser {
    return &RespParser{
        bufr: bufr,
    }
}

type RespObject interface {
    ToResp() string
}

type RespBulkString struct {
    data string
}

func (s RespBulkString) ToResp() string {
    return fmt.Sprintf("$%d\r\n%s\r\n", len(s.data), s.data)
}

type RespSimpleString struct {
    data string
}

func (s RespSimpleString) ToResp() string {
    return fmt.Sprintf("+%s\r\n", s.data)
}

type RespErrorString struct {
    errtype string
    message string
}

func (s RespErrorString) ToResp() string {
    return fmt.Sprintf("-%s %s\r\n", s.errtype, s.message)
}

type RespList struct {
    data []RespObject
}

func (s RespList) ToResp() string {
    respstrs := []string{fmt.Sprintf("*%d\r\n", len(s.data))}
    for _, item := range s.data {
        respstrs = append(respstrs, item.ToResp())
    }
    return strings.Join(respstrs, "")
}

type RespInteger struct {
    data int
}

func (s RespInteger) ToResp() string {
    return fmt.Sprintf(":%d\r\n", s.data)
}

func (p *RespParser) next() (RespObject, error) {
    begline, err := p.bufr.ReadString('\n')
    if err != nil {
        return nil, err
    }

    for {
        if begline[0] != '+' && begline[0] != '-' && begline[0] != '$' && begline[0] != '*' && begline[0] != ':' {
            begline, err = p.bufr.ReadString('\n')
            if err != nil {
                return nil, err
            }
            continue
        }
        break
    }

    switch begline[0] {
    case '+':
        return p.nextSimpleString(begline)
    case '-':
        return p.nextErrorString(begline)
    case '*':
        return p.nextList(begline)
    case '$':
        return p.nextBulkString(begline)
    case ':':
        return p.nextInteger(begline)
    }
    return nil, errors.New("Should not reach here")
}

func (p *RespParser) nextBulkString(begline string) (RespBulkString, error) {
    begline = strings.TrimRight(begline, "\r\n")
    begline = strings.TrimLeft(begline, "$")
    length, err := strconv.Atoi(begline)
    if err != nil {
        return RespBulkString{}, err
    }
    strline, err := p.bufr.ReadString('\n')
    if err != nil {
        return RespBulkString{}, err
    }
    for len(strline) < length+2 {
        s, err := p.bufr.ReadString('\n')
        if err != nil {
            return RespBulkString{}, err
        }
        strline += s
    }

    if len(strline) != length+2 {
        return RespBulkString{}, errors.New("Length of string not match")
    }

    return RespBulkString{
        data: strings.TrimRight(strline, "\r\n"),
    }, nil
}

func (p *RespParser) nextSimpleString(begline string) (RespSimpleString, error) {
    begline = strings.TrimRight(begline, "\r\n")
    begline = strings.TrimLeft(begline, "+")
    return RespSimpleString{
        data: begline,
    }, nil
}

func (p *RespParser) nextErrorString(begline string) (RespErrorString, error) {
    begline = strings.TrimRight(begline, "\r\n")
    begline = strings.TrimLeft(begline, "-")
    sp := strings.SplitN(begline, " ", 2)
    if len(sp) != 2 {
        return RespErrorString{}, errors.New("Invalid error string")
    }
    return RespErrorString{
        errtype: sp[0],
        message: sp[1],
    }, nil
}

func (p *RespParser) nextList(begline string) (RespList, error) {
    begline = strings.TrimRight(begline, "\r\n")
    begline = strings.TrimLeft(begline, "*")
    length, err := strconv.Atoi(begline)
    ret := RespList{}
    if err != nil {
        return ret, err
    }

    for i := 0; i < length; i++ {
        obj, err := p.next()
        if err != nil {
            return ret, err
        }
        ret.data = append(ret.data, obj)
    }

    return ret, nil
}

func (p *RespParser) nextInteger(begline string) (RespInteger, error) {
    begline = strings.TrimRight(begline, "\r\n")
    begline = strings.TrimLeft(begline, "*")
    integer, err := strconv.Atoi(begline)

    ret := RespInteger{}
    if err != nil {
        return ret, err
    }

    ret.data = integer
    return ret, nil
}
