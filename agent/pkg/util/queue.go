package util

import (
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"reflect"
)

type queue struct {
	items []codefresh.Environment
}

type Queue interface {
	Push(env codefresh.Environment) (bool, error)
	Notify()
}

var (
	existingQueue *queue
)

func Get() Queue {
	if existingQueue == nil {
		existingQueue = &queue{}
	}

	return existingQueue
}

func (q *queue) Push(env codefresh.Environment) (bool, error) {

	if q.items == nil {
		q.items = []codefresh.Environment{env}
	}

	if len(q.items) == 0 {
		q.items = append(q.items, env)
		return true, nil
	}

	lastEnv := q.items[len(q.items)-1]

	if reflect.DeepEqual(lastEnv, env) {
		return false, errors.New("queue already contains this element")
	}

	q.items = append(q.items, env)

	return true, nil
}

func (q *queue) Notify() {

	for _, item := range q.items {
		fmt.Println(item)
	}

}
