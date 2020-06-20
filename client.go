package tcpclient

import (
	"errors"
	"fmt"
	"github.com/JimmyTsai16/tcpclient/errorcode"
	logging "github.com/z9905080/gloger"
	"net"
	"strings"
	"time"
)

type ConnState int

const (
	Default ConnState = iota
	Open
	Close
)

type ReadFunc func([]byte, int)
type ConnStateFunc func(ConnState)
type ErrorFunc func(errorcode.ErrorCode, error)

type Client struct {
	addr       string
	conn       net.Conn
	ReadBuffer int64

	ReadHandler      ReadFunc
	ConnStateHandler ConnStateFunc
	ErrorHandler     ErrorFunc

	ReconnectDuration time.Duration
	reconnectTimes    int
	connectClosed     bool
	isAutoConnect     bool
}

func New() *Client {
	return &Client{
		addr:              "",
		conn:              nil,
		ReadBuffer:        128,
		ReadHandler:       func([]byte, int) {},
		ConnStateHandler:  func(ConnState) {},
		ErrorHandler:      func(errorcode.ErrorCode, error) {},
		ReconnectDuration: time.Second * 2,
		reconnectTimes:    0,

		connectClosed: true,
		isAutoConnect: false,
	}
}

func (c *Client) Connect(addr string, isAutoConnect bool) error {
	c.addr = addr
	c.isAutoConnect = isAutoConnect

	connErr := c.connect()
	if connErr != nil {
		if c.isAutoConnect {
			go c.reconnect()
		}
		return connErr
	}
	return nil
}

func (c *Client) connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = conn

	err = c.conn.SetDeadline(time.Time{})
	if err != nil {
		return err
	}

	c.connectClosed = false
	c.ConnStateHandler(Open)
	go c.readHandler()
	return nil
}

func (c *Client) reconnect() {
	c.reconnectTimes = 0
	for {
		c.reconnectTimes++
		logging.Info("reconnecting")
		if connErr := c.connect(); connErr == nil {
			logging.Info(fmt.Sprintf("reconnect successful, times: %d", c.reconnectTimes))
			return
		}
		logging.Error(fmt.Sprintf("reconnect fail, times: %d", c.reconnectTimes))
		time.Sleep(c.ReconnectDuration)
	}
}

// write bytes, error return -1 or return how many byte write
func (c *Client) WriteByte(b []byte) (int, error) {
	if c.connectClosed {
		return 0, fmt.Errorf("connection is closed")
	}
	return c.conn.Write(b)
}

func (c *Client) readHandler() {
	for {
		var b = make([]byte, c.ReadBuffer)
		n, err := c.conn.Read(b)
		if err != nil {
			if c.IsClosed() {
				logging.Error("connection have been closed")
				return
			}
			logging.Error("connection read error:", err)

			// default error
			errorCode := errorcode.ConnectionUnknownError
			if strings.Contains(err.Error(), "closed") {
				errorCode = errorcode.ConnectionClosed
			} else if strings.Contains(err.Error(), "unreachable") {
				errorCode = errorcode.ConnectionClosed
			}
			c.connectClosed = true

			closeErr := c.Close()
			if closeErr != nil {
				logging.Error("close connection error: closeErr:", closeErr)
			}

			//go c.ReadHandler(nil, -1)
			go c.ConnStateHandler(Close)
			go c.ErrorHandler(errorCode, errors.New(fmt.Sprint("Read error: ", err)))

			go c.reconnect()
			return
		}
		c.ReadHandler(b, n)
	}
}

func (c *Client) GetReconnectTimes() int {
	return c.reconnectTimes
}

func (c *Client) SetAutoReconnect(onOff bool) {
	c.isAutoConnect = onOff
}

func (c *Client) IsClosed() bool {
	return c.connectClosed
}

func (c *Client) Close() error {
	err := c.conn.Close()
	return err
}
