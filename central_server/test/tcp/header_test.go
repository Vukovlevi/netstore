package tcp_test

import (
	"encoding/binary"
	"log"
	"testing"

	"github.com/vukovlevi/netstore/central_server/tcp"
)

func TestCreateHeaderFromBuffer(t *testing.T) {
    testLength := 6

    headerBuffer := []byte{2}
    headerBuffer = binary.BigEndian.AppendUint32(headerBuffer, uint32(testLength))

    testWrongVersion(headerBuffer, testLength)
    testCorrectVersion(headerBuffer, testLength)
}

func TestCreateHeaderForPayload(t *testing.T) {
    payload := []byte{6, 32, 78, 90, 87, 0x4E}
    header := tcp.CreateHeaderForPayload(payload)

    if header[0] != tcp.VERSION {
        log.Fatalf("version sent by the server did not match expectations, expected: %d, got: %d", tcp.VERSION, header[0])
    }

    size := binary.BigEndian.Uint32(header[1:])
    if uint32(len(payload)) != size {
        log.Fatalf("header length does not match the length of actual payload, expected length: %d, actual length: %d", len(payload), size)
    }
}

func testWrongVersion(headerBuffer []byte, testLength int) {
    testVersion := 2
    headerBuffer[0] = byte(testVersion)

    header := testCreation(headerBuffer, testVersion, testLength)
    err := header.ValidateHeader()
    if err == nil {
        log.Fatal("expected to get an error")
    }

    if err != tcp.ErrVersionMismatch {
        log.Fatalf("expected version mismatch error, got: %s", err.Error())
    }
}

func testCorrectVersion(headerBuffer []byte, testLength int) {
    testVersion := 1
    headerBuffer[0] = byte(testVersion)

    header := testCreation(headerBuffer, testVersion, testLength)
    err := header.ValidateHeader()
    if err != nil {
        log.Fatalf("expected not to get an error")
    }
}

func testCreation(headerBuffer []byte, testVersion, testLength int) *tcp.TcpHeader {
    header := tcp.CreateHeaderFromBuffer(headerBuffer)
    if header == nil {
        log.Fatal("header cannot be nil")
    }

    if header.Version != byte(testVersion) {
        log.Fatal("version was not parsed correctly")
    }

    if header.MsgLen != uint32(testLength) {
        log.Fatalf("size of header is not correct, expected: %d, got: %d", 6, header.MsgLen)
    }

    return header
}
