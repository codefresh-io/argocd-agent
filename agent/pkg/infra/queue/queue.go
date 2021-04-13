package queue

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sync"
)

// ItemQueue the queue of Items
type ItemQueue struct {
	items []*unstructured.Unstructured
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
	s.items = make([]*unstructured.Unstructured, 0)
	return s
}

// Enqueue adds an Item to the end of the queue
func (s *ItemQueue) Enqueue(t *unstructured.Unstructured) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Dequeue removes an Item from the start of the queue
func (s *ItemQueue) Dequeue() *unstructured.Unstructured {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	s.lock.Unlock()
	return item
}

// Front returns the item next in the queue, without removing it
func (s *ItemQueue) Front() *unstructured.Unstructured {
	s.lock.RLock()
	item := s.items[0]
	s.lock.RUnlock()
	return item
}

// IsEmpty returns true if the queue is empty
func (s *ItemQueue) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of Items in the queue
func (s *ItemQueue) Size() int {
	return len(s.items)
}
