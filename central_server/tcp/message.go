package tcp

import (
	"errors"
	"os"
	"slices"
)

const (
    MSG_TYPE_AUTHENTICATION = byte(1)
    MSG_TYPE_SEARCH = byte(2)
    MSG_TYPE_ANSWER = byte(3)
    MSG_TYPE_CLIENT_SEARCH = byte(3)
    MSG_TYPE_CLIENT_ANSWER = byte(5)
    MSG_TYPE_ERROR = byte(6)

    MSG_EOF = byte(0x4E)
)

var (
    AuthenticationError = errors.New("the psk sent by the client did not match the one in server configuration")
)

type Message interface {
    ToMessageBytes() []byte
}

type TcpMessage struct {
    MessageType byte
    Content []byte
}

func (t *TcpMessage) ToMessageBytes() []byte {
    payload := []byte{t.MessageType}
    payload = append(payload, t.Content...)
    payload = append(payload, MSG_EOF)

    header := CreateHeaderForPayload(payload)

    message := slices.Concat(header, payload)
    return message
}

func CreateTcpMessageFromPayload(payload []byte) *TcpMessage {
    msgType := payload[0]
    content := payload[1:len(payload) - 1]
    return &TcpMessage{MessageType: msgType, Content: content}
}

type AuthenticationMessage struct {
    *TcpMessage
}

func (a *AuthenticationMessage) Authenticate() error {
    psk := os.Getenv("PSK")
    if psk != string(a.Content) {
        return AuthenticationError
    }

    return nil
}

type SearchMessage struct {
    *TcpMessage
}

type AnswerMessage struct {
    *TcpMessage
}

type ClientSearchMessage struct {
    *TcpMessage
}

type ClientAnswerMessage struct {
    *TcpMessage
}

type ErrorMessage struct {
    *TcpMessage
    Msg string
}

func CreateErrorMessage(msg string) *ErrorMessage {
    return &ErrorMessage{TcpMessage: &TcpMessage{MessageType: MSG_TYPE_ERROR}, Msg: msg}
}

func (e *ErrorMessage) ToMessageBytes() []byte {
    e.Content = []byte(e.Msg)
    return e.TcpMessage.ToMessageBytes()
}
