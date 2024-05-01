package lexigo

import (
	"bufio"
	"errors"
	"net"

	"github.com/vincer2040/lexiGo/internal/protocol"
	"github.com/vincer2040/lexiGo/pkg/lexitypes"
)

type Client struct {
	ip        string
	connected bool
	conn      net.Conn
}

func New(ip string) *Client {
	return &Client{
		ip:        ip,
		connected: false,
	}
}

func (c *Client) Close() {
	if !c.connected {
		return
	}
	c.conn.Close()
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.ip)
	if err != nil {
		return err
	}
	c.connected = true
	c.conn = conn
	return nil
}

func (c *Client) Ping() (string, error) {
	buf := protocol.NewBuilder().
		AddSimpleString("PING")
	err := c.write(buf)
	if err != nil {
		return "", err
	}
	data, err := c.read()
	if err != nil {
		return "", err
	}
	if data.DataType == lexitypes.Error {
		return "", errors.New(string(data.Data.(lexitypes.LexiString)))
	}
	if data.DataType != lexitypes.String {
		return "", errors.New("unexpected data type from server")
	}
	return string(data.Data.(lexitypes.LexiString)), nil
}

func (c *Client) Auth(username, password string) (string, error) {
	buf := protocol.NewBuilder().
		AddArray(3).
		AddSimpleString("AUTH").
		AddBulkString(username).
		AddBulkString(password)
	err := c.write(buf)
	if err != nil {
		return "", err
	}
	data, err := c.read()
	if err != nil {
		return "", err
	}
	if data.DataType == lexitypes.Error {
		return "", errors.New(string(data.Data.(lexitypes.LexiString)))
	}
	if data.DataType != lexitypes.String {
		return "", errors.New("unexpected data type from server")
	}
	return string(data.Data.(lexitypes.LexiString)), nil
}

func (c *Client) Set(key, value string) (string, error) {
	buf := protocol.NewBuilder().
		AddArray(3).
		AddSimpleString("SET").
		AddBulkString(key).
		AddBulkString(value)
	err := c.write(buf)
	if err != nil {
		return "", err
	}
	data, err := c.read()
	if err != nil {
		return "", err
	}
	if data.DataType == lexitypes.Error {
		return "", errors.New(string(data.Data.(lexitypes.LexiString)))
	}
	if data.DataType != lexitypes.String {
		return "", errors.New("unexpected data type from server")
	}
	return string(data.Data.(lexitypes.LexiString)), nil
}

func (c *Client) Get(key string) (string, error) {
	buf := protocol.NewBuilder().
		AddArray(2).
		AddSimpleString("GET").
		AddBulkString(key)
	err := c.write(buf)
	if err != nil {
		return "", err
	}
	data, err := c.read()
	if err != nil {
		return "", err
	}
	if data.DataType == lexitypes.Error {
		return "", errors.New(string(data.Data.(lexitypes.LexiString)))
	}
	if data.DataType == lexitypes.Null {
		return "", nil
	}
	if data.DataType != lexitypes.String {
		return "", errors.New("unexpected data type from server")
	}
	return string(data.Data.(lexitypes.LexiString)), nil
}

func (c *Client) Del(key string) (int64, error) {
	buf := protocol.NewBuilder().
		AddArray(2).
		AddSimpleString("DEL").
		AddBulkString(key)
	err := c.write(buf)
	if err != nil {
		return 0, err
	}
	data, err := c.read()
	if err != nil {
		return 0, err
	}
	if data.DataType == lexitypes.Error {
		return 0, errors.New(string(data.Data.(lexitypes.LexiString)))
	}
	if data.DataType != lexitypes.Int {
		return 0, errors.New("unexpected data type from server")
	}
	return int64(data.Data.(lexitypes.LexiInt)), nil
}

func (c *Client) Keys() ([]string, error) {
	buf := protocol.NewBuilder().
		AddSimpleString("KEYS")
	_, err := c.conn.Write(buf)
	if err != nil {
		return nil, err
	}
	data, err := c.read()
	if err != nil {
		return nil, err
	}
	if data.DataType == lexitypes.Error {
		return nil, errors.New(string(data.Data.(lexitypes.LexiString)))
	}
	if data.DataType != lexitypes.Array {
		return nil, errors.New("unexpected data type from server")
	}
	arr := data.Data.(lexitypes.LexiArray)
	res := make([]string, len(arr))
	for i, s := range arr {
		if s.DataType != lexitypes.String {
			return nil, errors.New("unexpected data type from server")
		}
		str := string(s.Data.(lexitypes.LexiString))
		res[i] = str
	}
	return res, nil
}

func (c *Client) Type(key string) (string, error) {
	buf := protocol.NewBuilder().
		AddArray(2).
		AddSimpleString("TYPE").
		AddBulkString(key)
	err := c.write(buf)
	if err != nil {
		return "", err
	}
	data, err := c.read()
	if err != nil {
		return "", err
	}
	if data.DataType == lexitypes.Error {
		return "", errors.New(string(data.Data.(lexitypes.LexiString)))
	}
	if data.DataType != lexitypes.String {
		return "", errors.New("unexpected data type from server")
	}
	return string(data.Data.(lexitypes.LexiString)), nil
}

func (c *Client) Incr(key string) (int64, error) {
	buf := protocol.NewBuilder().
		AddArray(2).
		AddSimpleString("INCR").
		AddBulkString(key)
	err := c.write(buf)
	if err != nil {
		return 0, err
	}
	data, err := c.read()
	if err != nil {
		return 0, err
	}
	if data.DataType == lexitypes.Error {
		return 0, errors.New(string(data.Data.(lexitypes.LexiString)))
	}
	if data.DataType != lexitypes.Int {
		return 0, errors.New("unexpected data type from server")
	}
	return int64(data.Data.(lexitypes.LexiInt)), nil
}

func (c *Client) write(buf []byte) error {
	length := len(buf)
	for length > 0 {
		n, err := c.conn.Write(buf)
		if err != nil {
			return err
		}
		buf = buf[n:]
		length -= n
	}
	return nil
}

func (c *Client) read() (*lexitypes.LexiType, error) {
	reader := protocol.NewReader(bufio.NewReader(c.conn))
	return reader.ReadReply()
}
