package queue

import "sync"

const (
    STATUS_CAN_SEARCH = 1
    STATUS_ANSWERING = 2
)

type SearchRequestNode struct {
    Next *SearchRequestNode
    SearchParam []byte
    AnswerChan chan []byte
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
        q.mutex.Unlock()
        return
    }

    curr := q.Dequeue()
    if curr == nil {
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

    q.Process()
}
