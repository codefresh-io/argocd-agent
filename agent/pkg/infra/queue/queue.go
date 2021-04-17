package queue

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"sync"
)

// ItemQueue the queue of Items
type ItemQueue struct {
	items map[string]*service.ApplicationWrapper
	lock  sync.RWMutex
}

var queueInstance = ItemQueue{
	items: nil,
	lock:  sync.RWMutex{},
}

var q *ItemQueue

func GetInstance() *ItemQueue {
	if q == nil {
		q = queueInstance.New()
	}
	return q
}

// New creates a new ItemQueue
func (s *ItemQueue) New() *ItemQueue {
	s.items = make(map[string]*service.ApplicationWrapper, 0)
	return s
}

// Enqueue adds an Item to the end of the queue
func (s *ItemQueue) Enqueue(t *service.ApplicationWrapper) {
	s.lock.Lock()
	key := fmt.Sprintf("%s.%v", t.Application.Status.OperationState.SyncResult.Revision, t.HistoryId)
	s.items[key] = t
	s.lock.Unlock()
}

// Dequeue removes an Item from the start of the queue
func (s *ItemQueue) Dequeue() *service.ApplicationWrapper {
	s.lock.Lock()
	for key, element := range s.items {
		if element != nil {
			delete(s.items, key)
			s.lock.Unlock()
			return element
		} else {
			delete(s.items, key)
		}
	}
	s.lock.Unlock()
	return nil
}

// Size returns the number of Items in the queue
func (s *ItemQueue) Size() int {
	return len(s.items)
}
