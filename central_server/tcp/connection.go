package tcp

import (
	"errors"
	"io"
	"log/slog"
	"net"

	"github.com/google/uuid"
)

const (
    VERSION_ERROR_MESSAGE = "not valid version"
    INVALID_MSG_TYPE_ERROR_MESSAGE = "not valid message type"
    AUTHENTICATION_ERROR_MESSAGE = "invalid credentials"
)

type Connection struct {
    Id uuid.UUID
    Conn net.Conn
    SearchRequestChan chan *SearchMessage
    AnswerChan chan *AnswerMessage
    CurrentAnswerId string
    IsAuthenticated bool
    ConnChan chan *Connection
    ReturnError error
}

func CreateConnection(conn net.Conn, searchRequestChan chan *SearchMessage, connChan chan *Connection) *Connection {
    return &Connection{
        Id: uuid.New(),
        Conn: conn,
        SearchRequestChan: searchRequestChan,
        ConnChan: connChan,
    }
}

func (c *Connection) ReadLoop() {
    defer func() {
        c.Conn.Close()
        if c.ReturnError == io.EOF {
            c.ConnChan <- c
        } else {
            slog.Error("connection forcefully closed by server", "error", c.ReturnError)
        }
    }()
    for {
        header, err := c.ReadHeader()
        if err != nil {
            if err == io.EOF {
                c.ReturnError = err
                return
            }
            continue
        }

        if err = header.ValidateHeader(); err != nil {
            c.ReadPayload(header.MsgLen)
            slog.Error("header not valid", "err", err)
            if err == ErrVersionMismatch {
                sendErr := c.SendErrorMessage(VERSION_ERROR_MESSAGE) //TODO: hungarian error message
                if sendErr != nil {
                    slog.Error("there was an error sending error message", "error", sendErr)
                }
            }
            continue
        }

        tcpMessage := c.ReadPayload(header.MsgLen)
        if tcpMessage == nil {
            continue
        }

        if !c.IsAuthenticated && tcpMessage.MessageType != MSG_TYPE_AUTHENTICATION {
            c.ReturnError = errors.New("client tried to send message without authenticating first")
            return
        }

        if c.IsAuthenticated && tcpMessage.MessageType == MSG_TYPE_AUTHENTICATION {
            continue
        }

        c.HandleMessage(tcpMessage)
    }
}

func (c *Connection) ReadHeader() (*TcpHeader, error) {
    headerBuffer := make([]byte, HEADER_SIZE)
    n, err := c.Conn.Read(headerBuffer)
    if err != nil {
        slog.Error("could not read client message", "error", err)
        return nil, err
    }

    if n != HEADER_SIZE {
        slog.Error("message from client is too short for a header", "message", headerBuffer[:n])
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
        slog.Error("error reading message", "error", err)
        return nil
    }

    if n != len(buffer) || buffer[len(buffer) - 1] != MSG_EOF {
        slog.Error("the read message did not match the length it said would have")
        return nil
    }

    return CreateTcpMessageFromPayload(buffer)
}

func (c *Connection) HandleMessage(message *TcpMessage) {
    switch message.MessageType {
    case MSG_TYPE_AUTHENTICATION:
        c.Authenticate(message.ToAuthenticationMessage())
    case MSG_TYPE_SEARCH:
        c.EnqueueSearchRequest(message.ToSearchMessage())
    case MSG_TYPE_ANSWER:
        c.GiveAnswer(message.ToAnswerMessage())
    default:
        c.SendErrorMessage(INVALID_MSG_TYPE_ERROR_MESSAGE) //TODO: hungarian error message
    }
}

func (c *Connection) Authenticate(message *AuthenticationMessage) {
    if err := message.Authenticate(); err != nil {
        slog.Error("authentication failure for client", "id", c.Id.String(), "address", c.Conn.RemoteAddr().String(), "sent psk", string(message.Content))
        err = c.SendErrorMessage(AUTHENTICATION_ERROR_MESSAGE) //TODO: hungarian error message
        if err != nil {
            slog.Error("could not send error message", "error", err)
        }
        return
    }
    c.IsAuthenticated = true
    c.ConnChan <- c
    c.SendMessage(CreateAuthenticationSuccessMessage())
}

func (c *Connection) EnqueueSearchRequest(message *SearchMessage) {
    c.SearchRequestChan <- message
}

func (c *Connection) GiveAnswer(message *AnswerMessage) {
    defer func() {
        err := recover()
        if err != nil {
            slog.Error("there was an error while answering search request", "error", err)
        }
    }()

    if message.AnswerId != c.CurrentAnswerId {
        return
    }
    c.AnswerChan <- message
}

func (c *Connection) SendErrorMessage(msg string) error {
    message := CreateErrorMessage(msg)
    return c.SendMessage(message)
}

func (c *Connection) SendMessage(message Message) error {
    send := message.ToMessageBytes()
    n, err := c.Conn.Write(send)
    if err != nil {
        return err
    }

    if n != len(send) {
        slog.Error("could not send whole message", "expected to send", len(send), "actually sent", n)
        return errors.New("failed to send whole message")
    }

    return nil
}
