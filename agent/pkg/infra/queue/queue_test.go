package queue

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestItemQueue(t *testing.T) {

	m := make(map[string]interface{})
	m["k"] = "v"

	queue := GetInstance()
	queue.Enqueue(&unstructured.Unstructured{Object: m})

	size := queue.Size()
	if size != 1 {
		t.Error("Wrong size of queue")
	}

	itm := queue.Dequeue()

	if itm == nil {
		t.Error("We should be able retrieve item")
	}

	queue.Enqueue(&unstructured.Unstructured{Object: m})

	queue = queue.New()
	size = queue.Size()
	if size != 0 {
		t.Error("Wrong size of queue after create new one")
	}

}
