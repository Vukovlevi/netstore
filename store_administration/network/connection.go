package network

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"sync"
)

var (
	AuthenticationError = errors.New("Csatlakozás a központi szerverhez sikertelen, azonosítás hiba.")
)

type Connection struct {
	Conn net.Conn
	ServerAnswerChan chan Message
	mutex *sync.Mutex
}

func ConnectToCentralServer(ip, port string) (*Connection, error) {
	if ip == "" {
		ip = os.Getenv("CENTRAL_SERVER_IP")
	}
	if port == "" {
		port = os.Getenv("CENTRAL_SERVER_PORT")
	}
	slog.Debug("connection to central server", "ip", ip, "port", port)
	address := fmt.Sprintf("%s:%s", ip, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Connection{Conn: conn, mutex: new(sync.Mutex)}, nil
}

//Returns human-readable error
func (c *Connection) Authenticate(psk string) error {
	err := c.SendMessage(CreateAuthenticationMessage(psk))
	if err != nil {
		slog.Error("could not send authentication message to server", "error", err)
		return AuthenticationError
	}
	header, customErr, networkErr := c.ReadHeader()
	if networkErr != nil || customErr != nil {
		slog.Error("could not read header from server after authentication", "error", err)
		return AuthenticationError
	}

	tcpMessage := c.ReadPayload(header.MsgLen)
	if tcpMessage == nil {
		slog.Error("could not create tcp message from server message after authentication")
		return AuthenticationError
	}

	if tcpMessage.MessageType == MSG_TYPE_ERROR {
		errorMessage := tcpMessage.ToErrorMessage()
		return errors.New(errorMessage.Msg)
	}

	if tcpMessage.MessageType != MSG_TYPE_AUTHENTICATION_SUCCESS {
		slog.Error("got unexpected msg type after authentication", "expected error", MSG_TYPE_ERROR, "or auth success", MSG_TYPE_AUTHENTICATION_SUCCESS, "got", tcpMessage.MessageType)
		return AuthenticationError
	}

	return nil
}

func (c *Connection) ReadLoop() {
    defer func() {
        c.Conn.Close()
    }()

    for {
        header, customErr, networkErr := c.ReadHeader()
		if networkErr != nil {
			slog.Error("there was an error reading from connection", "error", networkErr)
			return
		}

		if customErr != nil {
			continue
		}

        slog.Debug("header reading complete", "header", header)

        tcpMessage := c.ReadPayload(header.MsgLen)
        if tcpMessage == nil {
            continue
        }

        c.HandleMessage(tcpMessage)
    }
}

func (c *Connection) ReadHeader() (*TcpHeader, error, error) {
    headerBuffer := make([]byte, HEADER_SIZE)
    n, err := c.Conn.Read(headerBuffer)
    if err != nil {
        slog.Error("could not read server message", "error", err)
        return nil, nil, err
    }

    if n != HEADER_SIZE {
        slog.Error("message from server does not have enough bytes for a header", "message", headerBuffer[:n])
        return nil, errors.New("header size mismatch"), nil
    }

    return CreateHeaderFromBuffer(headerBuffer), nil, nil
}

func (c *Connection) ReadPayload(length uint32) *TcpMessage {
    if length < 2 {
        return nil
    }

    buffer := make([]byte, length)
    n, err := io.ReadFull(c.Conn, buffer)
    if err != nil {
        slog.Error("error reading message from server", "error", err)
        return nil
    }

    if n != len(buffer) || buffer[len(buffer) - 1] != MSG_EOF {
        slog.Error("the read message from server did not match the length it said would have")
        return nil
    }

    return CreateTcpMessageFromPayload(buffer)
}

func (c *Connection) HandleMessage(message *TcpMessage) {
    slog.Debug("handling message", "type", message.MessageType, "content", message.Content)
    switch message.MessageType {
    case MSG_TYPE_CLIENT_SEARCH:
        go c.GetSearchResults(message.ToClientSearchMessage())
    case MSG_TYPE_CLIENT_ANSWER:
        c.GiveServerAnswer(message.ToClientAnswerMessage())
    case MSG_TYPE_ERROR:
		c.Conn.Close()
        errorMessage := message.ToErrorMessage()
		c.GiveServerAnswer(errorMessage)
    default:
    }
}

func (c *Connection) GetSearchResults(message *ClientSearchMessage) {
	searchResult, err := Manager.GetSearchResults(message.SearchParam)
	if err != nil && err != ErrNoErrorMessage {
		slog.Error("could not get search result for client search", "error", err)
	}
	slog.Debug("got results for search request", "answer id", message.AnswerId, "result", searchResult)
	answerMessage := CreateAnswerMessage(message.AnswerId, searchResult)
	c.SendMessage(answerMessage)
}

func (c *Connection) GiveServerAnswer(message Message) {
	defer func() {
		err := recover()
		if err != nil {
			slog.Error("could not send full answer back to api", "error", err)
		}
	}()

	c.ServerAnswerChan <- message
}

func (c *Connection) SendSearchRequest(message *SearchMessage, answerChan chan Message) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.ServerAnswerChan = answerChan
	c.SendMessage(message)
}

func (c *Connection) write(msg []byte) error {
    n, err := c.Conn.Write(msg)
    if err != nil {
        return err
    }

    if n != len(msg) {
        slog.Error("could not send whole message to server", "expected to send", len(msg), "actually sent", n)
        return errors.New("failed to send whole message to server")
    }

    return nil
}

func (c *Connection) SendMessage(message Message) error {
    send := message.ToMessageBytes()
    return c.write(send)
}