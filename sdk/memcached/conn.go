package memcached

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

var (
	delimiter = []byte("\r\n")

	storedResp   = []byte("STORED\r\n")
	endResp      = []byte("END\r\n")
	deletedResp  = []byte("DELETED\r\n")
	notFoundResp = []byte("NOT_FOUND\r\n")

	valueHeaderRE = regexp.MustCompile(`^(?m)VALUE [a-zA-Z0-9_]+ \d+ (\d+)( \d+){0,1}\r\n$`) // value in first `()` is length.
)

const (
	maxKeySize   = 255
	maxValueSize = 1024 * 1024
)

type Conn struct {
	mu sync.Mutex
	c  net.Conn
	rw *bufio.ReadWriter
}

func NewConnWithAddress(addr string) (*Conn, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial")
	}

	return NewConn(c), nil
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		rw: bufio.NewReadWriter(
			bufio.NewReaderSize(conn, maxKeySize+maxValueSize+1024),
			bufio.NewWriterSize(conn, maxKeySize+maxValueSize+1024),
		),
		c: conn,
	}
}

func (c *Conn) Close() error {
	return c.c.Close()
}

func (c *Conn) Set(key string, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.rw.WriteString(set(key, 0, 0, len(value)))
	if err != nil {
		return err
	}
	_, err = c.rw.Write(delimiter)
	if err != nil {
		return err
	}
	_, err = c.rw.Write(value)
	if err != nil {
		return err
	}
	_, err = c.rw.Write(delimiter)
	if err != nil {
		return err
	}
	err = c.rw.Flush()
	if err != nil {
		return err
	}

	resp, err := c.rw.ReadBytes('\n')
	if err != nil {
		return err
	}

	if !bytes.Equal(resp, storedResp) {
		return ErrNotStored
	}

	return nil
}

func (c *Conn) Get(key string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.rw.WriteString(get(key))
	if err != nil {
		return nil, err
	}
	_, err = c.rw.Write(delimiter)
	if err != nil {
		return nil, err
	}
	err = c.rw.Flush()
	if err != nil {
		return nil, err
	}

	header, err := c.rw.ReadBytes('\n')
	if err != nil {
		return nil, errors.Wrap(err, "failed to read bytes")
	}

	if bytes.Equal(header, endResp) {
		return nil, ErrNotFound
	}

	matches := valueHeaderRE.FindSubmatch(header)
	if len(matches) == 0 {
		return nil, ErrInvalidValueHeader
	}

	length, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		panic(err) // should have been handled with regex.
	}

	res := make([]byte, length+7) // 7 is for `\r\nEND\r\n`.
	_, err = io.ReadFull(c.rw, res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read value")
	}

	return res[:len(res)-7], nil
}

func (c *Conn) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.rw.WriteString(delete(key))
	if err != nil {
		return err
	}
	_, err = c.rw.Write(delimiter)
	if err != nil {
		return err
	}
	err = c.rw.Flush()
	if err != nil {
		return err
	}

	resp, err := c.rw.ReadBytes('\n')
	if err != nil {
		return err
	}

	switch {
	case bytes.Equal(resp, deletedResp):
		return nil
	case bytes.Equal(resp, notFoundResp):
		return ErrNotFound
	}

	return errors.Wrap(ErrUnknownResponse, "failed to delete")
}

func set(key string, meta, expiry, length int) string {
	return fmt.Sprintf("set %s %d %d %d", key, meta, expiry, length)
}

func get(key string) string {
	return "get " + key
}

func delete(key string) string {
	return "delete " + key
}
