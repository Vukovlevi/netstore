package client_test

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"

	"github.com/vukovlevi/netstore/central_server/tcp"
)

type TestAnswer struct {

}

type TestClient struct {
	Conn net.Conn
}

func (c *TestClient) ConnectToServer() {
	conn, err := net.Dial("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatalf("could not connect to server, error: %s", err.Error())
	}
	c.Conn = conn
}

func (c *TestClient) TestUnauthenticated() {
	message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_CLIENT_ANSWER, Content: []byte(os.Getenv("PSK"))}
	c.SendMessage(&message)

	_, err := c.Conn.Read(make([]byte, 1024))
	if err != io.EOF {
		log.Fatal("expected eof error, because not authentication message was sent")
	}
}

func (c *TestClient) TestAuthentication() {
	message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_AUTHENTICATION, Content: []byte(os.Getenv("PSK"))}
	c.SendMessage(&message)
	headerBuffer := make([]byte, tcp.HEADER_SIZE)
	n, err := c.Conn.Read(headerBuffer)
	if err != nil {
		log.Fatalf("error reading authentication response header, error: %s", err.Error())
	}
	if n != 5 {
		log.Fatalf("expected header length to be 5, got: %d", n)
	}
	length := binary.BigEndian.Uint32(headerBuffer[1:])
	buffer := make([]byte, length)
	_, err = io.ReadFull(c.Conn, buffer)
	if err != nil {
		log.Fatalf("error while reading authentication response payload, error: %s", err.Error())
	}
	if buffer[0] != tcp.MSG_TYPE_AUTHENTICATION_SUCCESS {
		log.Fatalf("expected authentication success message (%d), got: %d", tcp.MSG_TYPE_AUTHENTICATION_SUCCESS, buffer[0])
	}
}

func (c *TestClient) SendMessage(message tcp.Message) {
	_, err := c.Conn.Write(message.ToMessageBytes())
	if err != nil {
		log.Fatalf("could not write test message to connection, error: %s", err.Error())
	}
}