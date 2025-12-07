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

func (n *SearchRequestNode) Process(processChan chan *SearchRequestNode) {
    processChan <- n
}

type SearchRequestQueue struct {
    Head *SearchRequestNode
    Tail *SearchRequestNode
    Status int
    SearchRequestChan chan *SearchRequestNode
    ProcessChan chan *SearchRequestNode
    mutex *sync.Mutex
}

func NewSearchRequestQueue() *SearchRequestQueue {
    return &SearchRequestQueue{
        Status: STATUS_CAN_SEARCH,
        SearchRequestChan: make(chan *SearchRequestNode, 1),
        ProcessChan: make(chan *SearchRequestNode, 1),
        mutex: new(sync.Mutex),
    }
}

func (q *SearchRequestQueue) Enqueue(node *SearchRequestNode) {
    if q.Head == nil {
        q.Head = node
        q.Tail = q.Head
    } else {
        q.Tail.Next = node
        q.Tail = node
    }

    q.Process()
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
    defer q.mutex.Unlock()

    if q.Status != STATUS_CAN_SEARCH {
        return
    }

    curr := q.Dequeue()
    if curr == nil {
        return
    }

    q.Status = STATUS_ANSWERING
    curr.Process(q.ProcessChan)
}

func (q *SearchRequestQueue) FinishProcess() {
    q.mutex.Lock()
    q.Status = STATUS_CAN_SEARCH
    q.mutex.Unlock()
    q.Process()
}
