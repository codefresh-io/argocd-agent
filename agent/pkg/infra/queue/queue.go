package queue

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
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
	if t.HistoryId == -1 {
		logger.GetLogger().Infof("Ignore add item to queue, history Id is unknown")
		return
	}
	if t.Application.Status.OperationState.SyncResult.Revision == "" {
		logger.GetLogger().Infof("Ignore add item to queue, revision is empty")
		return
	}
	s.lock.Lock()
	logger.GetLogger().Infof("Add item to queue, revision %v, history %v", t.Application.Status.OperationState.SyncResult.Revision, t.HistoryId)
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
		}
	}
	s.lock.Unlock()
	return nil
}

// Size returns the number of Items in the queue
func (s *ItemQueue) Size() int {
	return len(s.items)
}
