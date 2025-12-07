package tcp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/vukovlevi/netstore/central_server/queue"
)

type Server struct {
    Listener    net.Listener
    Connections map[string]*Connection
    ConnChan    chan *Connection
    SearchRequestQueue *queue.SearchRequestQueue
}

func NewServer() *Server {
	ln, err := net.Listen("tcp", "0.0.0.0:42069") //TODO: read from config
	if err != nil {
		slog.Error("could not create listener for server", "error", err)
		panic("missing listener")
	}
	return &Server{
		Listener:    ln,
		Connections: make(map[string]*Connection),
		ConnChan:    make(chan *Connection, 1),
        SearchRequestQueue: queue.NewSearchRequestQueue(),
	}
}

func (s *Server) Start() {
    go s.HandleConnections()
    go s.ProcessSearchRequest()

    for {
        conn, err := s.Listener.Accept()
        if err != nil {
            slog.Error("could not accept connection", "error", err)
            continue
        }

        connection := CreateConnection(conn, s.SearchRequestQueue.SearchRequestChan, s.ConnChan) //TODO: searchRequestChan
        go connection.ReadLoop()
    }
}

func (s *Server) HandleConnections() {
    for c := range s.ConnChan {
        if _, ok := s.Connections[c.Id.String()]; !ok {
            s.Connections[c.Id.String()] = c
        } else {
            delete(s.Connections, c.Id.String())
        }
    }
}

func (s *Server) ProcessSearchRequest() {
    for searchRequest := range s.SearchRequestQueue.ProcessChan {
        fmt.Println(searchRequest) //TODO: processing logic here (sending client search messages the collecting answers)
        s.SearchRequestQueue.FinishProcess()
    }
}
