package client_test

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/vukovlevi/netstore/central_server/tcp"
)

type TestAnswer struct {
	Username string `json:"username"`
	Role string `json:"role"`
	Num int `json:"num"`
}

const (
    TEST_KEY = "test"
    TEST_VALUE = "search"
)

var (
    TEST_SEARCH = map[string]string{TEST_KEY: TEST_VALUE}
)

func (a *TestAnswer) ToMessageBytes() []byte {
	content, err := json.Marshal(a)
	if err != nil {
		log.Fatalf("could not marshal test answer, error: %s", err.Error())
	}

	message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_ANSWER, Content: content}
	return message.ToMessageBytes()
}

type TestClient struct {
	Conn net.Conn
	Answer *TestAnswer
}

func NewTestClient(username, role string, num int) *TestClient {
	return &TestClient{
		Answer: &TestAnswer{Username: username, Role: role, Num: num},
	}
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
	buffer := c.Read()
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
	buffer := c.Read()
	if buffer[0] != tcp.MSG_TYPE_AUTHENTICATION_SUCCESS {
		log.Fatalf("expected authentication success message (%d) on auth success, got: %d", tcp.MSG_TYPE_AUTHENTICATION_SUCCESS, buffer[0])
	}
	if buffer[len(buffer) - 1] != tcp.MSG_EOF {
		log.Fatalf("expected last byte to be EOF (%d) on auth success, got: %d", tcp.MSG_EOF, buffer[len(buffer) - 1])
	}
}

func (c *TestClient) SendMalformedMessage() {
    message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_AUTHENTICATION, Content: []byte{1, 2, 3}}
    data := message.ToMessageBytes()
    binary.BigEndian.PutUint32(data[1:], 3)
    _, err := c.Conn.Write(data)
    if err != nil {
        log.Fatalf("error while sending malformed message, error: %s", err.Error())
    }
}

func (c *TestClient) SendSearchRequest() {
    data, _ := json.Marshal(TEST_SEARCH)
    searchRequest := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_SEARCH, Content: data}
    c.SendMessage(&searchRequest)
}

func (c *TestClient) SendAnswer(delay time.Duration) {
    time.Sleep(delay)
    c.Answer.Num++
    data, _ := json.Marshal(c.Answer)
    answer := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_ANSWER, Content: data}
    c.SendMessage(&answer)
}

func (c *TestClient) ReadClientSearch() {
    clientSearch := c.Read()
    if clientSearch[0] != tcp.MSG_TYPE_CLIENT_SEARCH {
        log.Fatalf("expected a client search message (%d), got: %d", tcp.MSG_TYPE_CLIENT_SEARCH, clientSearch[0])
    }
	if clientSearch[len(clientSearch) - 1] != tcp.MSG_EOF {
		log.Fatalf("expected last byte to be EOF (%d) on client search, got: %d", tcp.MSG_EOF, clientSearch[len(clientSearch) - 1])
	}
    clientSearch = clientSearch[1:len(clientSearch) - 1]
    data := make(map[string]string)
    if err := json.Unmarshal(clientSearch, &data); err != nil {
        log.Fatalf("error while unmarshalling client search, error: %s", err.Error())
    }
    test, ok := data[TEST_KEY]
    if !ok {
        log.Fatalf("expected %s to be present in client search as key, but it was not", TEST_KEY)
    }
    if test != TEST_VALUE {
        log.Fatalf("expected %s to be value of %s in client search", test, TEST_VALUE)
    }
}

func (c *TestClient) ReadClientAnswer(isEmpty bool, username, role string, num int) {
    clientAnswer := c.Read()
    if clientAnswer[0] != tcp.MSG_TYPE_CLIENT_ANSWER {
        log.Fatalf("expected a client answer message (%d), got: %d", tcp.MSG_TYPE_CLIENT_ANSWER, clientAnswer[0])
    }
	if clientAnswer[len(clientAnswer) - 1] != tcp.MSG_EOF {
		log.Fatalf("expected last byte to be EOF (%d) on client answer, got: %d", tcp.MSG_EOF, clientAnswer[len(clientAnswer) - 1])
	}
    clientAnswer = clientAnswer[1:len(clientAnswer) - 1]
    array := make([]TestAnswer, 0)
    if err := json.Unmarshal(clientAnswer, &array); err != nil {
        log.Fatalf("error while unmarshalling client answer, error: %s", err.Error())
    }

    if !isEmpty && len(array) == 0 {
        log.Fatalf("not expected client answer to be empty, but it is")
    } else if isEmpty {
        if len(array) == 0 {
            return
        }
        log.Fatalf("expected client answer to be empty, instead it is: %v", array)
    }

    data := array[0]
    if username != data.Username || role != data.Role || num != data.Num {
        log.Fatalf("client answer is not what is expected: username (%s->%s), role (%s->%s), num (%d->%d)", username, data.Username, role, data.Role, num ,data.Num)
    }
}

func (c *TestClient) Read() []byte {
	headerBuffer := make([]byte, tcp.HEADER_SIZE)
	n, err := c.Conn.Read(headerBuffer)
	if err != nil && err != io.EOF {
		log.Fatalf("error while reading connection, error: %s", err.Error())
	}
	if n != tcp.HEADER_SIZE {
		log.Fatalf("could not read a header from message from server, expected: %d, got: %d", tcp.HEADER_SIZE, n)
	}
	length := binary.BigEndian.Uint32(headerBuffer[1:])
	buffer := make([]byte, length)
	_, err = io.ReadFull(c.Conn, buffer)
	if err != nil {
		log.Fatalf("error while reading authentication response payload on auth success, error: %s", err.Error())
	}
	return buffer
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
