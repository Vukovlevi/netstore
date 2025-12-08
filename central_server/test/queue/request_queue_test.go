package queue_test

import (
	"log"
	"testing"

	"github.com/vukovlevi/netstore/central_server/queue"
)

var (
    third *queue.SearchRequestNode
    fourth *queue.SearchRequestNode
    called = 0
)

func TestQueueEvents(t *testing.T) {
    q := queue.NewSearchRequestQueue(testHandleProcess)
    q.IsTesting = true

    first := &queue.SearchRequestNode{SearchParam: []byte{1, 2, 3, 4}}
    second := &queue.SearchRequestNode{SearchParam: []byte{43, 178, 89}}
    third = &queue.SearchRequestNode{SearchParam: []byte{3, 88, 82, 87, 92}}
    fourth = &queue.SearchRequestNode{SearchParam: []byte{34, 82, 91, 27}}

    if q.Dequeue() != nil {
        t.Fatal("expected dequeue to return nil")
    }

    q.Enqueue(first)
    if q.Head != first {
        log.Fatalf("head mismatch, expected: %v, got: %v at pos1", first, q.Head)
    }
    if q.Head != q.Tail {
        log.Fatal("expected head to equal tail at pos1")
    }

    q.Enqueue(second)
    if q.Head != first {
        log.Fatalf("head mismatch, expected: %v, got: %v at pos2", first, q.Head)
    }
    if q.Tail != second {
        log.Fatalf("tail mismatch, expected: %v, got : %v at pos2", second, q.Tail)
    }
    if q.Head.Next != second {
        log.Fatalf("head next mismatch, expected: %v, got: %v at pos2", second, q.Head.Next)
    }

    q.Enqueue(third)
    if q.Head != first {
        log.Fatalf("head mismatch, expected: %v, got: %v at pos3", first, q.Head)
    }
    if q.Tail != third {
        log.Fatalf("tail mismatch, expected: %v, got : %v at pos3", third, q.Tail)
    }
    if q.Head.Next != second {
        log.Fatalf("head next mismatch, expected: %v, got: %v at pos3", second, q.Head.Next)
    }
    if second.Next != third {
        log.Fatalf("second next mismatch, expected: %v, got: %v at pos3", third, second.Next)
    }

    dq := q.Dequeue()
    if q.Head != second {
        log.Fatalf("head mismatch, expected: %v, got: $%v at pos4", second, q.Head)
    }
    if q.Head.Next != third {
        log.Fatalf("head next mismatch, expected: %v, got: %v at pos4", third, q.Head.Next)
    }
    if q.Tail != third {
        log.Fatalf("tail mismatch, expected: %v, got : %v at pos4", third, q.Tail)
    }
    if dq != first {
        log.Fatalf("dq mismatch, expected: %v, got: %v at pos4", first, dq)
    }

    dq = q.Dequeue()
    if q.Head != third {
        log.Fatalf("head mismatch, expected: %v, got: $%v at pos5", third, q.Head)
    }
    if q.Head.Next != nil {
        log.Fatalf("head next mismatch, expected: %v, got: %v at pos5", nil, q.Head.Next)
    }
    if q.Tail != third {
        log.Fatalf("tail mismatch, expected: %v, got : %v at pos5", third, q.Tail)
    }
    if dq != second {
        log.Fatalf("dq mismatch, expected: %v, got: %v at pos5", first, dq)
    }

    q.Process()
    if q.Status != queue.STATUS_ANSWERING {
        log.Fatalf("status mismatch, expected: %d, got: %d at pos6", queue.STATUS_ANSWERING, q.Status)
    }
    if called != 1 {
        log.Fatalf("callback call count mismatch, expected: %d, got: %d at pos6", 1, called)
    }
    if q.Head != nil {
        log.Fatalf("head mismatch, expected nil, got: %v at pos6", q.Head)
    }
    if q.Head != q.Tail {
        log.Fatalf("head and tail mismatch, head: %v, tail: %v at pos6", q.Head, q.Tail)
    }

    q.Process()
    if called != 1 {
        log.Fatalf("callback call count mismatch, expected: %d, got: %d at pos7", 1, called)
    }
    if q.Status != queue.STATUS_ANSWERING {
        log.Fatalf("status mismatch, expected: %d, got: %d at pos7", queue.STATUS_ANSWERING, q.Status)
    }

    q.FinishProcess()
    if called != 1 {
        log.Fatalf("callback call count mismatch, expected: %d, got: %d at pos8", 1, called)
    }
    if q.Status != queue.STATUS_CAN_SEARCH {
        log.Fatalf("status mismatch, expected: %d, got: %d at pos8", queue.STATUS_CAN_SEARCH, q.Status)
    }

    q.Enqueue(fourth)
    q.Process()
    if called != 2 {
        log.Fatalf("callback call count mismatch, expected: %d, got: %d at pos9", 2, called)
    }
    if q.Status != queue.STATUS_ANSWERING {
        log.Fatalf("status mismatch, expected: %d, got: %d at pos9", queue.STATUS_ANSWERING, q.Status)
    }
    q.FinishProcess()
    if called != 2 {
        log.Fatalf("callback call count mismatch, expected: %d, got: %d at pos10", 2, called)
    }
    if q.Status != queue.STATUS_CAN_SEARCH {
        log.Fatalf("status mismatch, expected: %d, got: %d at pos10", queue.STATUS_CAN_SEARCH, q.Status)
    }
}

func testHandleProcess(srn *queue.SearchRequestNode) {
    called++
    switch called {
    case 1:
        if srn != third {
            log.Fatalf("srn to be processed mismatch, expected: %v, got: %v at cb1", third, srn)
        }
        if srn.Next != nil {
            log.Fatalf("expected srn next to be nil, it was: %v at cb1", srn.Next)
        }
    case 2:
        if srn != fourth {
            log.Fatalf("srn to be processed mismatch, expected: %v, got: %v at cb2", fourth, srn)
        }
        if srn.Next != nil {
            log.Fatalf("expected srn next to be nil, it was: %v at cb2", srn.Next)
        }
    default:
        return
    }
}
