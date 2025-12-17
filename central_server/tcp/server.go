package tcp

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vukovlevi/netstore/central_server/queue"
)

const (
    TIMEOUT_IN_SECONDS = 3
)

type Server struct {
    Listener    net.Listener
    Connections map[string]*Connection
    ConnChan    chan *Connection
    SearchRequestQueue *queue.SearchRequestQueue
    mutex *sync.RWMutex
}

func NewServer() *Server {
	ln, err := net.Listen("tcp", "0.0.0.0:42069") //TODO: read from config
	if err != nil {
		slog.Error("could not create listener for server", "error", err)
		panic("missing listener")
	}
    server := &Server{
		Listener:    ln,
		Connections: make(map[string]*Connection),
		ConnChan:    make(chan *Connection, 1),
	}
    server.SearchRequestQueue = queue.NewSearchRequestQueue(server.ProcessSearchRequest)
    return server
}

func (s *Server) Start() {
    go s.HandleConnections()
    go s.SearchRequestQueue.HandleSearchRequest()

    for {
        conn, err := s.Listener.Accept()
        if err != nil {
            slog.Error("could not accept connection", "error", err)
            continue
        }

        slog.Debug("new connection accepted", "connection", conn)
        connection := CreateConnection(conn, s.SearchRequestQueue.SearchRequestChan, s.ConnChan) //TODO: searchRequestChan
        go connection.ReadLoop()
    }
}

func (s *Server) HandleConnections() {
    for c := range s.ConnChan {
        s.mutex.Lock()
        if _, ok := s.Connections[c.Id.String()]; !ok {
            slog.Debug("new connection added to connection list", "connection", c)
            s.Connections[c.Id.String()] = c
        } else {
            slog.Debug("deleting connection from connection list", "connection", c)
            delete(s.Connections, c.Id.String())
        }
        s.mutex.Unlock()
    }
}

func (s *Server) ProcessSearchRequest(searchRequest *queue.SearchRequestNode) {
    clientSearchMessage := CreateClientSearchMessage(searchRequest.ClientId, uuid.New().String(), searchRequest.SearchParam)
    s.BroadCastSearchMessage(clientSearchMessage, searchRequest.FullAnswerChan)
    s.SearchRequestQueue.FinishProcess()
}

func (s *Server) BroadCastSearchMessage(searchMessage *ClientSearchMessage, fullAnswerChan chan []byte) {
    s.mutex.RLock()
    wg := new(sync.WaitGroup)
    for id, connection := range s.Connections {
        if id == searchMessage.ClientId {
            continue
        }

        if err := connection.SendClientSearch(searchMessage); err != nil {
            slog.Error("could not send client search message", "error", err)
            continue
        }
        slog.Debug("send search message to connection", "message", searchMessage, "connection", connection)

        wg.Add(1)
    }
    s.mutex.RUnlock()

    s.ListenForAnswers(searchMessage.SingleAnswerChan, fullAnswerChan, wg)
}

func (s *Server) ListenForAnswers(singleAnswerChan chan *AnswerMessage, fullAnswerChan chan []byte, wg *sync.WaitGroup) {
    ctx, done := context.WithTimeout(context.Background(), time.Second * TIMEOUT_IN_SECONDS)
    answers := make([]*AnswerMessage, 0)

    wgDoneChan := make(chan struct{})
    go func(wgDoneChan chan struct{}, wg *sync.WaitGroup) {
        wg.Wait()
        wgDoneChan <- struct{}{}
    }(wgDoneChan, wg)

    for {
        select {
        case singleAnswer := <- singleAnswerChan:
            slog.Debug("got single answer", "answer", singleAnswer)
            answers = append(answers, singleAnswer)
            wg.Done()
        case <- wgDoneChan:
        case <- ctx.Done():
            slog.Debug("waiting for clients to answer is done", "len of answers", len(answers))
            done()
            close(singleAnswerChan)
            s.CreateAndSendClientAnswer(answers, fullAnswerChan)
        }
    }
}

func (s *Server) CreateAndSendClientAnswer(singleAnswers []*AnswerMessage, fullAnswerChan chan []byte) {
    fullAnswer := make([]map[string]any, 0)
    for _, singleAnswer := range singleAnswers {
        answer := make(map[string]any)
        if err := json.Unmarshal(singleAnswer.Content, &answer); err != nil {
            slog.Error("could not unmarshal single answer", "error", err)
            continue
        }
        fullAnswer = append(fullAnswer, answer)
    }

    clientAnswerContent, err := json.Marshal(fullAnswer)
    if err != nil {
        slog.Error("could not marshal client answer content", "error", err)
        errMsg := CreateErrorMessage("error while encoding answer") //TODO: hungarian error message
        fullAnswerChan <- errMsg.ToMessageBytes()
        return
    }

    clientAnswerMessage := CreateClientAnswerMessage(clientAnswerContent)
    slog.Debug("created client answer", "message", clientAnswerMessage)
    fullAnswerChan <- clientAnswerMessage.ToMessageBytes()
}
