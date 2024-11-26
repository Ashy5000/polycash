package node_util

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

var Conn net.Conn

func EstablishConnection() net.Listener {
	fmt.Println("Preparing connection...")
	listener, err := net.Listen("unix", "/tmp/vm.sock")
	if err != nil {
		panic(err)
	}

	fmt.Println("Waiting for connection...")
	connLocal, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	Conn = connLocal

	fmt.Println("Done!")
	return listener
}

func SendString(conn net.Conn, message string) error {
	lenBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuf, uint32(len(message)))

	_, err := conn.Write(lenBuf)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func ReceiveString() (string, error) {
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(Conn, lenBuf)
	if err != nil {
		return "", err
	}
	length := binary.LittleEndian.Uint32(lenBuf)
	res := make([]byte, length)
	_, err = io.ReadFull(Conn, res)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func CloseConnection(listener net.Listener) {
	err := listener.Close()
	if err != nil {
		panic(err)
	}
	if err = os.Remove("/tmp/vm.sock"); err != nil {
		panic(err)
	}
	err = Conn.Close()
	if err != nil {
		panic(err)
	}
}
