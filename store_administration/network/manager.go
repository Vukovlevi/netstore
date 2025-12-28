package network

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

const (
	TIMEOUT_IN_SECONDS = 30
	STATUS_NOT_CONNECTED = 1
	STATUS_CAN_SEARCH = 2
	STATUS_WAITING_FOR_ANSWER = 3
)

type NetworkManager struct {
	Connection *Connection
	Status int
	mutex *sync.RWMutex
}

var Manager *NetworkManager

// Returns human-readable error
func NewNetworkManager(ip, port, psk string) error {
	Manager = &NetworkManager{Status: STATUS_NOT_CONNECTED}
	conn, err := ConnectToCentralServer(ip, port)
	if err != nil {
		slog.Error("could not connect to central server", "error", err)
		return errors.New("Csatlakozás a központi szerverhez sikertelen!")
	}
	err = conn.Authenticate(psk)
	if err != nil {
		slog.Error("authentication failure in manager")
		return err
	}
	Manager = &NetworkManager{Connection: conn, Status: STATUS_CAN_SEARCH, mutex: new(sync.RWMutex)}
	go Manager.StartReadLoop()
	return nil
}

func (n *NetworkManager) StartReadLoop() {
	n.Connection.ReadLoop()
}

func (n *NetworkManager) CanSearch() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.Status != STATUS_NOT_CONNECTED
}

//Returns user-readable error message
func (n *NetworkManager) SearchNetwork(searchParam []byte) ([]byte, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.Status = STATUS_WAITING_FOR_ANSWER

	answerChan := make(chan Message, 1)
	message := CreateSearchMessage(searchParam)
	n.Connection.SendSearchRequest(message, answerChan)
	ctx, done := context.WithTimeout(context.Background(), time.Second * TIMEOUT_IN_SECONDS)

	var err error = nil
	answer := []byte{}
	select {
	case message := <- answerChan:
		if errorMessage, ok := message.(*ErrorMessage); ok {
			err = errors.New(errorMessage.Msg)
		} else if clientAnswerMessage, ok := message.(*ClientAnswerMessage); ok {
			answer = clientAnswerMessage.Content
		} else {
			err = errors.New("Ismeretlen hiba miatt a keresési kérelem nem sikerült!")
		}
	case <- ctx.Done():
		err = errors.New("A keresési kérelem sikertelen, mert túllépte az időkorlátot!")
	}

	done()
	close(answerChan)
	n.Status = STATUS_CAN_SEARCH
	return answer, err
}

func (n *NetworkManager) GetSearchResults(searchParam []byte) []byte {
	//TODO: get products from dave api, get store details and open hours from DB, bundle into a JSON object, marshal it then return that as a result
	return []byte{}
}