package network

import (
	"errors"
	"io"
	"log/slog"
	"net"
)

var (
	AuthenticationError = errors.New("Csatlakozás a központi szerverhez sikertelen, azonosítás hiba.")
)

type Connection struct {
	Conn net.Conn
}

func ConnectToCentralServer() (*Connection, error) {
	conn, err := net.Dial("tcp", "localhost:42069") //TODO: read from config
	if err != nil {
		return nil, err
	}
	return &Connection{Conn: conn}, nil
}

//Returns human-readable error
func (c *Connection) Authenticate() error {
	defer func() {
		c.Conn.Close()
	}()

	err := c.SendMessage(CreateAuthenticationMessage())
	if err != nil {
		slog.Error("could not send authentication message to server", "error", err)
		return AuthenticationError
	}
	header, err := c.ReadHeader()
	if err != nil {
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
        header, err := c.ReadHeader()
        if err != nil {
            if err == io.EOF {
                return
            }
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

func (c *Connection) ReadHeader() (*TcpHeader, error) {
    headerBuffer := make([]byte, HEADER_SIZE)
    n, err := c.Conn.Read(headerBuffer)
    if err != nil {
        slog.Error("could not read server message", "error", err)
        return nil, err
    }

    if n != HEADER_SIZE {
        slog.Error("message from server does not have enough bytes for a header", "message", headerBuffer[:n])
        return nil, errors.New("header size mismatch")
    }

    return CreateHeaderFromBuffer(headerBuffer), nil
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
        //TODO
    case MSG_TYPE_CLIENT_ANSWER:
        //TODO
    case MSG_TYPE_ERROR:
        //TODO
    default:
    }
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