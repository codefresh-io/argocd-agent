package queue

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
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

	var env argoSdk.ArgoApplication

	util.Convert(unstructured.Unstructured{Object: m}, env)
	queue.Enqueue(&env)

	size := queue.Size()
	if size != 1 {
		t.Error("Wrong size of queue")
	}

	itm := queue.Dequeue()

	if itm == nil {
		t.Error("We should be able retrieve item")
	}

	queue.Enqueue(&env)

	queue = queue.New()
	size = queue.Size()
	if size != 0 {
		t.Error("Wrong size of queue after create new one")
	}

}