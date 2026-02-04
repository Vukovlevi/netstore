package tcp

import (
	"encoding/binary"
	"errors"
)

const (
    VERSION = byte(1)

    VERSION_SIZE = 1
    MESSAGE_LENGTH_SIZE = 4
    HEADER_SIZE = VERSION_SIZE + MESSAGE_LENGTH_SIZE
)

var (
    ErrVersionMismatch = errors.New("the version of client message does not match the version of server")
)

type TcpHeader struct {
    Version byte
    MsgLen uint32
}

func CreateHeaderFromBuffer(header []byte) *TcpHeader {
    return &TcpHeader{
        Version: header[0],
        MsgLen: binary.BigEndian.Uint32(header[1:]),
    }
}

func (h *TcpHeader) ValidateHeader() error {
    if h.Version != VERSION {
        return ErrVersionMismatch
    }
    return nil
}

func CreateHeaderForPayload(payload []byte) []byte {
    header := make([]byte, HEADER_SIZE)
    header[0] = VERSION

    msgLen:= len(payload)
    binary.BigEndian.PutUint32(header[1:], uint32(msgLen))

    return header
}
