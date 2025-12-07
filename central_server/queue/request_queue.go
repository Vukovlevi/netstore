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

func (n *SearchRequestNode) Process(processCallback func(*SearchRequestNode)) {
    processCallback(n)
}

type SearchRequestQueue struct {
    Head *SearchRequestNode
    Tail *SearchRequestNode
    Status int
    SearchRequestChan chan *SearchRequestNode
    ProcessCallback func(*SearchRequestNode)
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

    return curr
}

func (q *SearchRequestQueue) Process() {
    q.mutex.Lock()
    if q.Status != STATUS_CAN_SEARCH {
        return
    }

    curr := q.Dequeue()
    if curr == nil {
        return
    }

    q.Status = STATUS_ANSWERING
    q.mutex.Unlock()

    curr.Process(q.ProcessCallback)
}
