package rpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

const numOfDataLength = 8

// ReadMsg 协议: 包头: 8个字节为数据长度
func ReadMsg(conn net.Conn) ([]byte, error) {

	msgLenBytes := make([]byte, numOfDataLength)
	_, err := io.ReadFull(conn, msgLenBytes)
	defer func() {
		if msg := recover(); msg != nil {
			err = errors.New(fmt.Sprintf("%v", msg))
		}
	}()
	if err != nil {
		return nil, err
	}

	dataLen := binary.BigEndian.Uint64(msgLenBytes)
	bs := make([]byte, dataLen)
	_, err = io.ReadFull(conn, bs)
	return bs, err
}

func readN(conn net.Conn, len uint64) ([]byte, error) {
	res := make([]byte, len)
	var readNum uint64 = 0
	for {
		n, err := conn.Read(res[readNum:])
		if err != nil {
			return nil, err
		}
		readNum += uint64(n)
		if readNum >= len {
			return res, nil
		}
	}
}

// EncodeMsg 加上头部,  表示data长度
func EncodeMsg(data []byte) []byte {
	dataLen := len(data)
	res := make([]byte, dataLen+numOfDataLength)
	binary.BigEndian.PutUint64(res[:numOfDataLength], uint64(dataLen))
	copy(res[numOfDataLength:], data)
	return res
}

func WriteMsg(conn net.Conn, data []byte) error {

	bs := EncodeMsg([]byte(data))
	// write 不像read,  go规定 read可以不读满buffer,  并且不返回error,
	//  而write如果不把buffer中的数据都写进去,    必须要返回error.

	// Write must return a non-nil error if it returns n < len(p).
	// Write must not modify the slice data, even temporarily.
	_, err := conn.Write(bs)
	return err
}
