package network

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/vukovlevi/netstore/store_administration/model"
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

type SearchResult struct {
	OpenHours []model.OpenHour `json:"open_hours"`
	StoreDetail model.StoreDetail `json:"store_detail"`
	Products any `json:"products"`
}

var (
	ErrNoErrorMessage = errors.New("could not marshal error message")
)

// Returns human-readable error
func NewNetworkManager(ip, port, psk string) error {
	Manager = &NetworkManager{Status: STATUS_NOT_CONNECTED, mutex: new(sync.RWMutex)}
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

func (n *NetworkManager) IsConnected() bool {
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

func (n *NetworkManager) GetSearchResults(searchParam []byte) ([]byte, error) {
	errBytes, err := json.Marshal(map[string]string{"error": "A keresés sikertelen!"})
	if err != nil {
		slog.Error("could not marshal error message on search result error", "error", err)
		return []byte{}, ErrNoErrorMessage
	}

	searchResult := SearchResult{}
	openHours, err := model.GetOpenHours(true)
	if err != nil {
		return errBytes, err
	}
	searchResult.OpenHours = openHours

	storeDetail, err := model.GetStoreDetail()
	if err != nil {
		return errBytes, err
	}
	storeDetail.CentralServerAddress = ""
	storeDetail.CentralServerPort = 0
	storeDetail.StoreTypeId = 0
	searchResult.StoreDetail = storeDetail

	//TODO: external api call here
	if err = json.Unmarshal(searchParam, &searchResult.Products); err != nil {
		return errBytes, err
	}

	searchResultBytes, err := json.Marshal(searchResult)
	if err != nil {
		return errBytes, err
	}
	return searchResultBytes, nil
}