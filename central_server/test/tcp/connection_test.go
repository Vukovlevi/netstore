package tcp_test

import (
	"crypto/rand"
	"io"
	"log"
	"net"
	"sync"
	"testing"

	"github.com/vukovlevi/netstore/central_server/tcp"
)

func TestReading_ShortSuccess(t *testing.T) {
    client, server := net.Pipe()
    defer client.Close()
    defer server.Close()

    connection := tcp.Connection{Conn: server}
    content := []byte{2, 3, 8, 97}
    message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_SEARCH, Content: content}

    go func() {
        send := message.ToMessageBytes()
        n, err := client.Write(send)
        if err != nil {
            log.Fatalf("error writing to connection short success, error: %s", err.Error())
        }

        if n != len(send) {
            log.Fatalf("could not write full message short success, expected to write: %d, written: %d", len(send), n)
        }
    }()

    header := connection.ReadHeader()
    if header == nil {
        log.Fatalf("header should not be nil here")
    }

    err := header.ValidateHeader()
    if err != nil {
        log.Fatalf("header should be valid, error: %s", err.Error())
    }

    readMessage := connection.ReadPayload(header.MsgLen)
    if readMessage == nil {
        log.Fatalf("read message should not be nil here")
    }

    if readMessage.MessageType != message.MessageType {
        log.Fatalf("message type mismatch, expected: %d, got: %d", message.MessageType, readMessage.MessageType)
    }

    for i, item := range readMessage.Content {
        if item != content[i] {
            log.Fatalf("content mismatch at (%d->%d), expected: %d, got: %d", i, i, content[i], item)
        }
    }
}

func TestReading_VersionMismatch(t *testing.T) {
    client, server := net.Pipe()
    defer client.Close()
    defer server.Close()

    connection := tcp.Connection{Conn: server}
    content := []byte{2, 3, 8, 97}
    message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_SEARCH, Content: content}

    go func() {
        send := message.ToMessageBytes()
        send[0] = 2
        n, err := client.Write(send)
        if err != nil {
            log.Fatalf("error writing to connection version mismatch, error: %s", err.Error())
        }

        if n != len(send) {
            log.Fatalf("could not write full message version mismatch, expected to write: %d, written: %d", len(send), n)
        }
    }()

    header := connection.ReadHeader()
    if header == nil {
        log.Fatalf("header should not be nil here")
    }

    err := header.ValidateHeader()
    if err == nil {
        log.Fatal("header should be invalid")
    }

    connection.ReadPayload(header.MsgLen)

    go func() {
        err := connection.SendErrorMessage(tcp.VERSION_ERROR_MESSAGE)
        if err != nil {
            log.Fatalf("could not send error message, error: %s", err.Error())
        }
    }()

    buffer := make([]byte, 1024)
    n, err := client.Read(buffer)
    if err != nil {
        log.Fatalf("error while reading error message, error: %s", err.Error())
    }

    message = *tcp.CreateTcpMessageFromPayload(buffer[5:n])
    if message.MessageType != tcp.MSG_TYPE_ERROR {
        log.Fatalf("message type mismatch on receiving error, expected: %d, got: %d", tcp.MSG_TYPE_ERROR, message.MessageType)
    }

    if string(message.Content) != tcp.VERSION_ERROR_MESSAGE {
        log.Fatalf("expected to get version error, got: %s", string(message.Content))
    }
}

func TestReading_LongSuccess(t *testing.T) {
    client, server := net.Pipe()
    defer client.Close()
    defer server.Close()

    connection := tcp.Connection{Conn: server}
    content := make([]byte, 4000)
    io.ReadFull(rand.Reader, content)
    message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_SEARCH, Content: content}

    wg := new(sync.WaitGroup)
    wg.Add(2)

    go func(wg *sync.WaitGroup) {
        send := message.ToMessageBytes()
        n, err := client.Write(send[:3000])
        if err != nil {
            log.Fatalf("error writing to connection long success, error: %s", err.Error())
        }

        if n != 3000 {
            log.Fatalf("could not write full message long success, expected to write: %d, written: %d", 3000, n)
        }

        n, err = client.Write(send[3000:])
        if err != nil {
            log.Fatalf("error writing to connection, error: %s", err.Error())
        }

        if n != len(send) - 3000 {
            log.Fatalf("could not write full message, expected to write: %d, written: %d", len(send) - 3000, n)
        }

        wg.Done()
    }(wg)


    go func(wg *sync.WaitGroup) {
        header := connection.ReadHeader()
        if header == nil {
            log.Fatalf("header should not be nil here")
        }

        err := header.ValidateHeader()
        if err != nil {
            log.Fatalf("header should be valid, error: %s", err.Error())
        }

        if header.MsgLen != 4002 {
            log.Fatalf("msg len mismatch, expected: %d, got: %d", 4002, header.MsgLen)
        }

        readMessage := connection.ReadPayload(header.MsgLen)
        if readMessage == nil {
            log.Fatalf("read message should not be nil here")
        }

        if readMessage.MessageType != message.MessageType {
            log.Fatalf("message type mismatch, expected: %d, got: %d", message.MessageType, readMessage.MessageType)
        }

        for i, item := range readMessage.Content {
            if item != content[i] {
                log.Fatalf("content mismatch at (%d->%d), expected: %d, got: %d", i, i, content[i], item)
            }
        }

        wg.Done()
    }(wg)

    wg.Wait()
}
