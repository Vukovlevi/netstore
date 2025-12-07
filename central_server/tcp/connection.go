package tcp

import (
	"errors"
	"io"
	"log/slog"
	"net"

	"github.com/google/uuid"
)

type Connection struct {
    Id uuid.UUID
    Conn io.ReadWriteCloser
    SearchRequestChan chan<- SearchMessage
    AnswerChan chan<- AnswerMessage
    CurrentAnswerId string
}

func CreateConnection(conn net.Conn, searchRequestChan chan<- SearchMessage) *Connection {
    return &Connection{
        Id: uuid.New(),
        Conn: conn,
        SearchRequestChan: searchRequestChan,
    }
}

func (c *Connection) ReadLoop() {
    for {
        header := c.ReadHeader()
        if header == nil {
            continue
        }

        if err := header.ValidateHeader(); err != nil {
            slog.Error("header not valid", "err", err)
            if err == ErrVersionMismatch {
                c.SendErrorMessage("not valid version") //TODO: hungarian error message
            }
            continue
        }

        tcpMessage := c.ReadPayload(header.GetLength())
        if tcpMessage == nil {
            continue
        }
    }
}

func (c *Connection) ReadHeader() *TcpHeader {
    headerBuffer := make([]byte, HEADER_SIZE)
    n, err := c.Conn.Read(headerBuffer)
    if err != nil {
        slog.Error("could not read client message", "error", err)
        return nil
    }

    if n != HEADER_SIZE {
        slog.Error("message from client is too short for a header", "message", headerBuffer[:n])
        return nil
    }

    return CreateHeaderFromBuffer(headerBuffer)
}

func (c *Connection) ReadPayload(length int) *TcpMessage {
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
