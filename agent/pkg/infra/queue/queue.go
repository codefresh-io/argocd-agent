package queue

import (
	"fmt"
	"sync"

	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
)

type AppEventsQueue struct {
	items []*service.ApplicationWrapper
	lock  sync.RWMutex
}

var queue = &AppEventsQueue{
	items: make([]*service.ApplicationWrapper, 0),
	lock:  sync.RWMutex{},
}

func GetAppQueue() *AppEventsQueue {
	return queue
}

// Enqueue adds an Item to the end of the queue
func (q *AppEventsQueue) Enqueue(event *service.ApplicationWrapper) {
	if event.HistoryId == -1 {
		logger.GetLogger().Infof("Ignore add item to queue, history Id is unknown")
		return
	}
	if event.Application.Status.OperationState.SyncResult.Revision == "" {
		logger.GetLogger().Infof("Ignore add item to queue, revision is empty")
		return
	}
	q.lock.Lock()
	defer q.lock.Unlock()

	logger.GetLogger().Infof("Add item to queue, revision %v, history %v", event.Application.Status.OperationState.SyncResult.Revision, event.HistoryId)

	q.items = append(q.items, event)
}

// Dequeue removes an Item from the start of the queue
func (q *AppEventsQueue) Dequeue() *service.ApplicationWrapper {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) == 0 {
		return nil
	}

	dequeuedItem := q.items[0]
	q.items = q.items[1:]

	fmt.Printf("dequeued item: %v, history id: %v", dequeuedItem.Application.Status.Sync.Revision, dequeuedItem.HistoryId)
	return dequeuedItem
}

// Size returns the number of Items in the queue
func (s *AppEventsQueue) Size() int {
	return len(s.items)
}

func (s *AppEventsQueue) Purge() {
	s.items = make([]*service.ApplicationWrapper, 0)
}
