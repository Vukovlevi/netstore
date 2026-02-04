package network

import (
	"os"
	"slices"
)

const (
    MSG_TYPE_AUTHENTICATION = byte(1)
    MSG_TYPE_SEARCH = byte(2)
    MSG_TYPE_ANSWER = byte(3)
    MSG_TYPE_CLIENT_SEARCH = byte(4)
    MSG_TYPE_CLIENT_ANSWER = byte(5)
    MSG_TYPE_ERROR = byte(6)
    MSG_TYPE_AUTHENTICATION_SUCCESS = byte(7)

    MSG_EOF = byte(0x4E)

    UUID_LENGTH = 36
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

func (m *TcpMessage) ToClientSearchMessage() *ClientSearchMessage {
    answerId := string(m.Content[:UUID_LENGTH])
    m.Content = m.Content[UUID_LENGTH:]
    return &ClientSearchMessage{
        TcpMessage: &TcpMessage{MessageType: MSG_TYPE_CLIENT_SEARCH},
        AnswerId: answerId,
        SearchParam: m.Content,
    }
}

func (m *TcpMessage) ToClientAnswerMessage() *ClientAnswerMessage {
    return &ClientAnswerMessage{TcpMessage: &TcpMessage{MessageType: MSG_TYPE_CLIENT_ANSWER, Content: m.Content}}
}

func (m *TcpMessage) ToErrorMessage() *ErrorMessage {
    return &ErrorMessage{TcpMessage: &TcpMessage{MessageType: MSG_TYPE_ERROR}, Msg: string(m.Content)}
}

type AuthenticationMessage struct {
    *TcpMessage
    Psk string
}

func (a *AuthenticationMessage) ToMessageBytes() []byte {
    a.Content = []byte(a.Psk)
    return a.TcpMessage.ToMessageBytes()
}

func CreateAuthenticationMessage(psk string) *AuthenticationMessage {
    if psk == "" {
        psk = os.Getenv("PSK")
    }
    return &AuthenticationMessage{TcpMessage: &TcpMessage{MessageType: MSG_TYPE_AUTHENTICATION}, Psk: psk}
}

type SearchMessage struct {
    *TcpMessage
}

type ClientAnswerMessage struct {
    *TcpMessage
}

func CreateSearchMessage(searchParam []byte) *SearchMessage {
    return &SearchMessage{TcpMessage: &TcpMessage{MessageType: MSG_TYPE_SEARCH, Content: searchParam}}
}

type AnswerMessage struct {
    *TcpMessage
    AnswerId string
    Answer []byte
}

func CreateAnswerMessage(answerId string, answer []byte) *AnswerMessage {
    return &AnswerMessage{TcpMessage: &TcpMessage{MessageType: MSG_TYPE_ANSWER}, AnswerId: answerId, Answer: answer}
}

func (a *AnswerMessage) ToMessageBytes() []byte {
    a.Content = slices.Concat([]byte(a.AnswerId), a.Answer)
    return a.TcpMessage.ToMessageBytes()
}

type ClientSearchMessage struct {
    *TcpMessage
    AnswerId string
    SearchParam []byte
}

type ErrorMessage struct {
    *TcpMessage
    Msg string
}