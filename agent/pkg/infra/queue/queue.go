package queue

import (
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"sync"
)

// ItemQueue the queue of Items
type ItemQueue struct {
	items []*argoSdk.ArgoApplication
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
	s.items = make([]*argoSdk.ArgoApplication, 0)
	return s
}

// Enqueue adds an Item to the end of the queue
func (s *ItemQueue) Enqueue(t *argoSdk.ArgoApplication) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an Item from the start of the queue
func (s *ItemQueue) Dequeue() *argoSdk.ArgoApplication {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	s.lock.Unlock()
	return item
}

// Size returns the number of Items in the queue
func (s *ItemQueue) Size() int {
	return len(s.items)
}
