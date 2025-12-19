package tcp

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/vukovlevi/netstore/central_server/queue"
)

const (
    VERSION_ERROR_MESSAGE = "not valid version"
    INVALID_MSG_TYPE_ERROR_MESSAGE = "not valid message type"
    AUTHENTICATION_ERROR_MESSAGE = "invalid credentials"
)

type Connection struct {
    Id uuid.UUID
    Conn net.Conn
    SearchRequestChan chan *queue.SearchRequestNode
    AnswerChan chan *AnswerMessage
    CurrentAnswerId string
    IsAuthenticated bool
    ConnChan chan *Connection
    ReturnError error
    mutex *sync.RWMutex
}

func CreateConnection(conn net.Conn, searchRequestChan chan *queue.SearchRequestNode, connChan chan *Connection) *Connection {
    return &Connection{
        Id: uuid.New(),
        Conn: conn,
        SearchRequestChan: searchRequestChan,
        ConnChan: connChan,
        mutex: new(sync.RWMutex),
    }
}

func (c *Connection) ReadLoop() {
    defer func() {
        c.Conn.Close()
        if c.ReturnError == io.EOF && c.IsAuthenticated {
            c.ConnChan <- c
        } else if c.ReturnError != io.EOF {
            slog.Error("connection forcefully closed by server", "client id", c.Id.String(), "error", c.ReturnError)
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
        slog.Debug("header reading complete", "client id", c.Id.String(), "header", header)

        if err = header.ValidateHeader(); err != nil {
            c.ReadPayload(header.MsgLen)
            slog.Error("header not valid", "client id", c.Id.String(), "err", err)
            if err == ErrVersionMismatch {
                sendErr := c.SendErrorMessage(VERSION_ERROR_MESSAGE) //TODO: hungarian error message
                if sendErr != nil {
                    slog.Error("there was an error sending error message", "client id", c.Id, "error", sendErr)
                }
            }
            continue
        }

        tcpMessage := c.ReadPayload(header.MsgLen)
        if tcpMessage == nil {
            continue
        }

        if !c.IsAuthenticated && tcpMessage.MessageType != MSG_TYPE_AUTHENTICATION {
            slog.Debug("client tried to send messages without authentication", "client id", c.Id.String())
            c.ReturnError = errors.New("client tried to send message without authenticating first")
            return
        }

        if c.IsAuthenticated && tcpMessage.MessageType == MSG_TYPE_AUTHENTICATION {
            slog.Debug("authenticated client tried to authenticate again", "client id", c.Id.String())
            continue
        }

        c.HandleMessage(tcpMessage)
    }
}

func (c *Connection) ReadHeader() (*TcpHeader, error) {
    headerBuffer := make([]byte, HEADER_SIZE)
    n, err := c.Conn.Read(headerBuffer)
    if err != nil {
        slog.Error("could not read client message", "client id", c.Id.String(), "error", err)
        return nil, err
    }

    if n != HEADER_SIZE {
        slog.Error("message from client is too short for a header", "client id", c.Id.String(), "message", headerBuffer[:n])
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
        slog.Error("error reading message", "client id", c.Id.String(), "error", err)
        return nil
    }

    if n != len(buffer) || buffer[len(buffer) - 1] != MSG_EOF {
        slog.Error("the read message did not match the length it said would have", "client id", c.Id.String())
        return nil
    }

    return CreateTcpMessageFromPayload(buffer)
}

func (c *Connection) HandleMessage(message *TcpMessage) {
    slog.Debug("handling message", "client id", c.Id.String(), "type", message.MessageType, "content", message.Content)
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
    slog.Debug("authentication done", "client", c.Id.String())
    c.IsAuthenticated = true
    c.ConnChan <- c
    c.SendMessage(CreateAuthenticationSuccessMessage())
}

func (c *Connection) EnqueueSearchRequest(message *SearchMessage) {
    node := &queue.SearchRequestNode{
        SearchParam: message.Content,
        FullAnswerChan: message.FullAnswerChan,
        ClientId: c.Id.String(),
    }
    c.SearchRequestChan <- node
    go c.WaitForAnswer(message.FullAnswerChan)
}

func (c *Connection) WaitForAnswer(answerChan chan []byte) {
    answer := <- answerChan
    slog.Debug("got client answer from chanel", "client id", c.Id.String(), "answer", answer)
    if err := c.write(answer); err != nil {
        slog.Error("could not send client answer message", "error", err, "client id", c.Id.String(), "content", answer)
    }
    close(answerChan)
    slog.Debug("client answer sent successfully", "client id", c.Id.String())
}

func (c *Connection) GiveAnswer(message *AnswerMessage) {
    defer func() {
        err := recover()
        if err != nil {
            slog.Error("there was an error while answering search request", "client id", c.Id.String(), "error", err)
        }
    }()

    c.mutex.RLock()
    defer c.mutex.RUnlock()
    if message.AnswerId != c.CurrentAnswerId {
        slog.Debug("got an answer for the wrong search request", "client id", c.Id.String(), "current answer id", c.CurrentAnswerId, "got answer id", message.AnswerId)
        return
    }
    c.AnswerChan <- message
}

func (c *Connection) SendErrorMessage(msg string) error {
    message := CreateErrorMessage(msg)
    return c.SendMessage(message)
}

func (c *Connection) write(msg []byte) error {
    n, err := c.Conn.Write(msg)
    if err != nil {
        return err
    }

    if n != len(msg) {
        slog.Error("could not send whole message", "client id", c.Id.String(), "expected to send", len(msg), "actually sent", n)
        return errors.New("failed to send whole message")
    }

    return nil
}

func (c *Connection) SendMessage(message Message) error {
    send := message.ToMessageBytes()
    return c.write(send)
}

func (c *Connection) SendClientSearch(message *ClientSearchMessage) error {
    c.mutex.Lock()
    c.AnswerChan = message.SingleAnswerChan
    c.CurrentAnswerId = message.AnswerId
    c.mutex.Unlock()
    return c.SendMessage(message)
}
