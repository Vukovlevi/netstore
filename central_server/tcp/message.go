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
    MSG_TYPE_AUTHENTICATION_SUCCESS = byte(7)

    MSG_EOF = byte(0x4E)

    UUID_LENGTH = 36
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

func (m *TcpMessage) ToMessageBytes() []byte {
    payload := []byte{m.MessageType}
    payload = append(payload, m.Content...)
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

func (m *TcpMessage) ToAuthenticationMessage() *AuthenticationMessage {
    return &AuthenticationMessage{TcpMessage: m, Psk: string(m.Content)}
}

func (m *TcpMessage) ToSearchMessage() *SearchMessage {
    return &SearchMessage{TcpMessage: m, AnswerChan: make(chan []*AnswerMessage, 1)}
}

func (m *TcpMessage) ToAnswerMessage() *AnswerMessage {
    answerId := string(m.Content[:UUID_LENGTH])
    m.Content = m.Content[UUID_LENGTH:]
    return &AnswerMessage{TcpMessage: m, AnswerId: answerId}
}

type AuthenticationMessage struct {
    *TcpMessage
    Psk string
}

func (a *AuthenticationMessage) Authenticate() error {
    psk := os.Getenv("PSK")
    if psk != a.Psk {
        return AuthenticationError
    }
    return nil
}

type SearchMessage struct {
    *TcpMessage
    AnswerChan chan []*AnswerMessage
}

type AnswerMessage struct {
    *TcpMessage
    AnswerId string
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

type AuthenticationSuccessMessage struct {
    *TcpMessage
}

func CreateAuthenticationSuccessMessage() *AuthenticationSuccessMessage {
    return &AuthenticationSuccessMessage{TcpMessage: &TcpMessage{MessageType: MSG_TYPE_AUTHENTICATION_SUCCESS}}
}
