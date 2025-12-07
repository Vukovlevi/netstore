package tcp_test

import (
	"crypto/rand"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
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

    message = *tcp.CreateTcpMessageFromPayload(buffer[tcp.HEADER_SIZE:n])
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

func TestHandleMessage(t *testing.T) {
    godotenv.Load()

    client, server := net.Pipe()
    defer client.Close()
    defer server.Close()

    connection := tcp.Connection{Conn: server}
    content := []byte{2, 3, 8, 97}
    message := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_AUTHENTICATION, Content: content}

    go sendMessageToServer(client, &message)
    testAuthenticationFailure(&connection, client)

    message.Content = []byte(os.Getenv("PSK"))
    go sendMessageToServer(client, &message)
    testAuthenticationSuccess(&connection, client)

    message.MessageType = tcp.MSG_TYPE_SEARCH
    go sendMessageToServer(client, &message)
    testEnqueueSearchRequest(&connection)

    message.MessageType = tcp.MSG_TYPE_ANSWER
    id := uuid.New().String()
    connection.CurrentAnswerId = id
    message.Content = []byte(id)
    message.Content = append(message.Content, []byte{2, 48, 29, 179}...)
    go sendMessageToServer(client, &message)
    testGiveAnswer(&connection)

    message.MessageType = tcp.MSG_TYPE_ANSWER
    connection.CurrentAnswerId = ""
    go sendMessageToServer(client, &message)
    testGiveInvalidAnswer(&connection)

    message.MessageType = tcp.MSG_TYPE_CLIENT_ANSWER
    go sendMessageToServer(client, &message)
    testInvalidMsgType(&connection, client)
}

func sendMessageToServer(client net.Conn, message tcp.Message) {
    send := message.ToMessageBytes()
    n, err := client.Write(send)
    if err != nil {
        log.Fatalf("error writing to connection short success, error: %s", err.Error())
    }

    if n != len(send) {
        log.Fatalf("could not write full message short success, expected to write: %d, written: %d", len(send), n)
    }
}

func testAuthenticationFailure(connection *tcp.Connection, client net.Conn) {
    header := connection.ReadHeader()
    readMessage := connection.ReadPayload(header.MsgLen)
    go connection.HandleMessage(readMessage)

    buffer := make([]byte, 1024)
    n, err := client.Read(buffer)
    if err != nil {
        log.Fatalf("error reading authentication failure, error: %s", err.Error())
    }

    if tcp.AUTHENTICATION_ERROR_MESSAGE != string(buffer[tcp.HEADER_SIZE + 1:n - 1]) {
        log.Fatalf("authentication failure mismatch, expected: %s, got: %s", tcp.AUTHENTICATION_ERROR_MESSAGE, string(buffer[tcp.HEADER_SIZE + 1:n - 1]))
    }
}

func testAuthenticationSuccess(connection *tcp.Connection, client net.Conn) {
    connection.NewConnChan = make(chan *tcp.Connection, 1)
    header := connection.ReadHeader()
    readMessage := connection.ReadPayload(header.MsgLen)
    go connection.HandleMessage(readMessage)

    buffer := make([]byte, 1024)
    n, err := client.Read(buffer)
    if err != nil {
        log.Fatalf("error reading authentication failure, error: %s", err.Error())
    }

    if n != tcp.HEADER_SIZE + 2 { //+2 -> msg type + eof => full payload for this type of message
        log.Fatalf("authentication success message length mismatch, expected: %d, got: %d", tcp.HEADER_SIZE + 2, n)
    }

    if buffer[tcp.HEADER_SIZE] != tcp.MSG_TYPE_AUTHENTICATION_SUCCESS {
        log.Fatalf("authentication message type mismatch, expected: %d, got: %d", tcp.MSG_TYPE_AUTHENTICATION_SUCCESS, buffer[tcp.HEADER_SIZE])
    }

    c := <- connection.NewConnChan
    if c != connection {
        log.Fatalf("did not get connection on new con chang")
    }

    close(connection.NewConnChan)
}

func testEnqueueSearchRequest(connection *tcp.Connection) {
    connection.SearchRequestChan = make(chan *tcp.SearchMessage, 1)

    header := connection.ReadHeader()
    readMessage := connection.ReadPayload(header.MsgLen)
    connection.HandleMessage(readMessage)

    searchMessage := <- connection.SearchRequestChan
    if searchMessage == nil {
        log.Fatalf("expected a search message on this chanel")
    }

    close(connection.SearchRequestChan)
}

func testGiveAnswer(connection *tcp.Connection) {
    connection.AnswerChan = make(chan *tcp.AnswerMessage, 1)

    header := connection.ReadHeader()
    readMessage := connection.ReadPayload(header.MsgLen)
    connection.HandleMessage(readMessage)

    answer := <- connection.AnswerChan
    if answer == nil {
        log.Fatalf("expected an answer on this chanel")
    }

    close(connection.AnswerChan)
}

func testGiveInvalidAnswer(connection *tcp.Connection) {
    header := connection.ReadHeader()
    readMessage := connection.ReadPayload(header.MsgLen)
    connection.HandleMessage(readMessage)
    //works as a test, because in case of invalid answer id, it should not send anything on the channel -> no error
    //if it is trying to send to closed chanel -> test doesnt fail because of recover, but error message will be shown
}

func testInvalidMsgType(connection *tcp.Connection, client net.Conn) {
    header := connection.ReadHeader()
    readMessage := connection.ReadPayload(header.MsgLen)
    go connection.HandleMessage(readMessage)

    buffer := make([]byte, 1024)
    n, err := client.Read(buffer)
    if err != nil {
        log.Fatalf("error reading authentication failure, error: %s", err.Error())
    }

    if buffer[tcp.HEADER_SIZE] != tcp.MSG_TYPE_ERROR {
        log.Fatalf("excpected to get error message on invalid msg type, got: %d", buffer[tcp.HEADER_SIZE])
    }

    if tcp.INVALID_MSG_TYPE_ERROR_MESSAGE != string(buffer[tcp.HEADER_SIZE + 1:n-1]) {
        log.Fatalf("error message mismatch on invalid msg type, expected: %s, got: %s", tcp.INVALID_MSG_TYPE_ERROR_MESSAGE, string(buffer[tcp.HEADER_SIZE + 1:n - 1]))
    }
}
