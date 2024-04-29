package client

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
	_, err := c.conn.Write(buf)
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
	_, err := c.conn.Write(buf)
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
	_, err := c.conn.Write(buf)
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
	_, err := c.conn.Write(buf)
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
	_, err := c.conn.Write(buf)
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
	for length < 0 {
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
