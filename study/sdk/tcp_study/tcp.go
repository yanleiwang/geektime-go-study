package tcp_study

import (
	"encoding/binary"
	"net"
)

const numOfDataLength = 8

// ReadMsg 协议: 包头: 8个字节为数据长度
func ReadMsg(conn net.Conn) ([]byte, error) {

	lenBs, err := readN(conn, numOfDataLength)
	if err != nil {
		return nil, err
	}

	return readN(conn, binary.BigEndian.Uint64(lenBs))

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

func EncodeMsg(data []byte) []byte {
	dataLen := len(data)
	res := make([]byte, dataLen+numOfDataLength)
	binary.BigEndian.PutUint64(res[:numOfDataLength], uint64(dataLen))
	copy(res[numOfDataLength:], data)
	return res
}

func WriteMsg(conn net.Conn, data []byte) error {
	bs := EncodeMsg([]byte(data))
	writedNum := 0
	for {
		n, err := conn.Write(bs[writedNum:])
		if err != nil {
			return err
		}
		writedNum += n
		if writedNum >= len(data) {
			return nil
		}
	}

}
