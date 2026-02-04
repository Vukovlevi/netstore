package tcp_test

import (
	"log"
	"testing"

	"github.com/vukovlevi/netstore/central_server/tcp"
)

func TestCreateTcpMessageFromPayload(t *testing.T) {
    messageType := byte(1)
    payload := []byte{messageType, 4, 68, 89, 21, 0x4E}
    tcpMessage := tcp.CreateTcpMessageFromPayload(payload)

    if tcpMessage.MessageType != messageType {
        log.Fatalf("message type did not match expectations, expected: %d, got: %d", messageType, tcpMessage.MessageType)
    }

    if len(tcpMessage.Content) != len(payload) - 2 {
        log.Fatalf("content length did not match expectations, expected: %d, got: %d", len(payload) - 2, len(tcpMessage.Content))
    }

    for i, item := range tcpMessage.Content {
        if item != payload[i + 1] {
            log.Fatalf("payload and content mismatch at (%d->%d), expected content to be: %d, got: %d", i + 1, i, payload[i + 1], item)
        }
    }
}

func TestToMessageBytes(t *testing.T) {
    content := []byte{34, 87, 92, 28}
    tcpMessage := tcp.TcpMessage{MessageType: tcp.MSG_TYPE_CLIENT_ANSWER, Content: content}
    testToMessageBytes(tcp.MSG_TYPE_CLIENT_ANSWER, content, tcpMessage.ToMessageBytes())
}

func TestErrorToMessageBytes(t *testing.T) {
    errMsg := "there was an error"
    errorMsg := tcp.CreateErrorMessage(errMsg)

    if errorMsg.MessageType != tcp.MSG_TYPE_ERROR {
        log.Fatalf("error message didnt have error msg type, got type: %d", errorMsg.MessageType)
    }

    if errorMsg.Msg != errMsg {
        log.Fatalf("error message didnt match expectation, expected: %s, got: %s", errMsg, errorMsg.Msg)
    }

    testToMessageBytes(tcp.MSG_TYPE_ERROR, []byte(errMsg), errorMsg.ToMessageBytes())
}

func testToMessageBytes(msgType byte, content, send []byte) {
    payload := []byte{msgType}
    payload = append(payload, content...)
    payload = append(payload, tcp.MSG_EOF)

    header := tcp.CreateHeaderForPayload(payload)
    sendHeader := send[:5]

    for i, item := range sendHeader {
        if item != header[i] {
            log.Fatalf("message bytes header mismatch at (%d->%d), expected: %d, got: %d", i, i, header[i], item)
        }
    }

    sendPayload := send[5:]
    if len(payload) != len(sendPayload) {
        log.Fatalf("payload lengths not matching, expected: %d, got: %d", len(payload), len(sendPayload))
    }

    for i, item := range sendPayload {
        if item != payload[i] {
            log.Fatalf("message bytes payload mismatch at (%d->%d), expected: %d, got: %d", i, i, payload[i], item)
        }
    }
}
