package queue

import (
	"log/slog"
	"sync"
)

const (
    STATUS_CAN_SEARCH = 1
    STATUS_ANSWERING = 2
)

type SearchRequestNode struct {
    Next *SearchRequestNode
    SearchParam []byte
    FullAnswerChan chan []byte
    ClientId string
}

type SearchRequestQueue struct {
    Head *SearchRequestNode
    Tail *SearchRequestNode
    Status int
    SearchRequestChan chan *SearchRequestNode
    ProcessCallback func(*SearchRequestNode)
    IsTesting bool
    mutex *sync.Mutex
}

func NewSearchRequestQueue(processCallback func(*SearchRequestNode)) *SearchRequestQueue {
    return &SearchRequestQueue{
        Status: STATUS_CAN_SEARCH,
        SearchRequestChan: make(chan *SearchRequestNode, 1),
        ProcessCallback: processCallback,
        mutex: new(sync.Mutex),
    }
}

func (q *SearchRequestQueue) HandleSearchRequest() {
    for searchRequest := range q.SearchRequestChan {
        q.Enqueue(searchRequest)
        slog.Debug("got search request into queue", "client id", searchRequest.ClientId, "search param", string(searchRequest.SearchParam))
    }
}

func (q *SearchRequestQueue) Enqueue(node *SearchRequestNode) {
    q.mutex.Lock()
    if q.Head == nil {
        q.Head = node
        q.Tail = q.Head
    } else {
        q.Tail.Next = node
        q.Tail = node
    }
    q.mutex.Unlock()

    if q.IsTesting {
        return
    }
    go q.Process()
}

func (q *SearchRequestQueue) Dequeue() *SearchRequestNode {
    if q.Head == nil {
        return nil
    }

    curr := q.Head
    q.Head = q.Head.Next
    if q.Head == nil {
        q.Tail = nil
    }

    curr.Next = nil
    return curr
}

func (q *SearchRequestQueue) Process() {
    q.mutex.Lock()
    if q.Status != STATUS_CAN_SEARCH {
        slog.Debug("not starting processing of search request, because there is one already processing")
        q.mutex.Unlock()
        return
    }

    curr := q.Dequeue()
    if curr == nil {
        slog.Debug("not starting processing of search request, because the queue is empty")
        q.mutex.Unlock()
        return
    }

    q.Status = STATUS_ANSWERING
    q.mutex.Unlock()

    q.ProcessCallback(curr)
}

func (q *SearchRequestQueue) FinishProcess() {
    q.mutex.Lock()
    q.Status = STATUS_CAN_SEARCH
    q.mutex.Unlock()

    slog.Debug("finished processing of search request")
    q.Process()
}
