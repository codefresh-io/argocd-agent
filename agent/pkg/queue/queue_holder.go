package queue

import "sync"

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
