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
	conn, err := net.Dial("tcp", "127.0.0.1:42069") //TODO: read from config
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

func (c *TestClient) TestAuthenticationFailure() {
	message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_AUTHENTICATION, Content: []byte(os.Getenv("PSK") + "-")}
	c.SendMessage(&message)
	headerBuffer := make([]byte, tcp.HEADER_SIZE)
	n, err := c.Conn.Read(headerBuffer)
	if err != nil {
		log.Fatalf("error reading authentication response header on auth fail, error: %s", err.Error())
	}
	if n != 5 {
		log.Fatalf("expected header length to be 5 on auth fail, got: %d", n)
	}
	length := binary.BigEndian.Uint32(headerBuffer[1:])
	buffer := make([]byte, length)
	_, err = io.ReadFull(c.Conn, buffer)
	if err != nil {
		log.Fatalf("error while reading authentication response payload on auth fail, error: %s", err.Error())
	}
	if buffer[0] != tcp.MSG_TYPE_ERROR {
		log.Fatalf("expected error message (%d) on auth fail, got: %d", tcp.MSG_TYPE_AUTHENTICATION_SUCCESS, buffer[0])
	}
	if buffer[len(buffer) - 1] != tcp.MSG_EOF {
		log.Fatalf("expected last byte to be EOF (%d) on auth fail, got: %d", tcp.MSG_EOF, buffer[len(buffer) - 1])
	}
	if string(buffer[1:len(buffer) - 1]) != tcp.AUTHENTICATION_ERROR_MESSAGE {
		log.Fatalf("expected authentication error message (%s), got: %s", tcp.AUTHENTICATION_ERROR_MESSAGE, string(buffer[1:len(buffer) - 1]))
	}
}

func (c *TestClient) TestAuthenticationSuccess() {
	message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_AUTHENTICATION, Content: []byte(os.Getenv("PSK"))}
	c.SendMessage(&message)
	headerBuffer := make([]byte, tcp.HEADER_SIZE)
	n, err := c.Conn.Read(headerBuffer)
	if err != nil {
		log.Fatalf("error reading authentication response header on auth successs, error: %s", err.Error())
	}
	if n != 5 {
		log.Fatalf("expected header length to be 5 on auth success, got: %d", n)
	}
	length := binary.BigEndian.Uint32(headerBuffer[1:])
	buffer := make([]byte, length)
	_, err = io.ReadFull(c.Conn, buffer)
	if err != nil {
		log.Fatalf("error while reading authentication response payload on auth success, error: %s", err.Error())
	}
	if buffer[0] != tcp.MSG_TYPE_AUTHENTICATION_SUCCESS {
		log.Fatalf("expected authentication success message (%d) on auth success, got: %d", tcp.MSG_TYPE_AUTHENTICATION_SUCCESS, buffer[0])
	}
	if buffer[len(buffer) - 1] != tcp.MSG_EOF {
		log.Fatalf("expected last byte to be EOF (%d) on auth success, got: %d", tcp.MSG_EOF, buffer[len(buffer) - 1])
	}
}

func (c *TestClient) SendMessage(message tcp.Message) {
	_, err := c.Conn.Write(message.ToMessageBytes())
	if err != nil {
		log.Fatalf("could not write test message to connection, error: %s", err.Error())
	}
}

func (c *TestClient) Close() {
	c.Conn.Close()
}