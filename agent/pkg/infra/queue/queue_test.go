package queue

import (
	"testing"

	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"github.com/stretchr/testify/assert"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestQueueSmoke(t *testing.T) {
	queue.Purge()

	aw1 := newApplicationWrapper(1)
	aw2 := newApplicationWrapper(2)
	aw3 := newApplicationWrapper(3)
	aw4 := newApplicationWrapper(4)

	queue.Enqueue(aw1)
	queue.Enqueue(aw2)
	queue.Enqueue(aw3)
	queue.Enqueue(aw4)

	assert.Equal(t, 4, queue.Size(), "Queue size should be 3.")

	item1 := queue.Dequeue()
	item2 := queue.Dequeue()
	item3 := queue.Dequeue()
	item4 := queue.Dequeue()

	assert.Equal(t, int64(1), item1.HistoryId, "Queue order is distorted.")
	assert.Equal(t, int64(2), item2.HistoryId, "Queue order is distorted.")
	assert.Equal(t, int64(3), item3.HistoryId, "Queue order is distorted.")
	assert.Equal(t, int64(4), item4.HistoryId, "Queue order is distorted.")

	assert.Equal(t, 0, queue.Size(), "Queue size should be 0.")
}

func TestItemQueueWithWrongHistoryId(t *testing.T) {
	queue.Purge()

	queue.Enqueue(newApplicationWrapper(-1))

	assert.Equal(t, 0, queue.Size(), "Queue size should be 0.")
}

func TestItemQueueWithWrongRevision(t *testing.T) {
	queue.Purge()

	// Application.Status.OperationState.SyncResult.Revision is nil
	queue.Enqueue(&service.ApplicationWrapper{
		Application: argoSdk.ArgoApplication{},
		HistoryId:   1,
	})

	assert.Equal(t, 0, queue.Size(), "Queue size should be 0.")
}

func TestItemQueueDequeueEmptyState(t *testing.T) {
	queue.Purge()

	assert.Equal(t, 0, queue.Size(), "Queue size should be 0.")

	item := queue.Dequeue()

	assert.Nil(t, item, "Should return nil when queue is empty")
}

func newApplicationWrapper(historyId int64) *service.ApplicationWrapper {
	aw := service.ApplicationWrapper{
		Application: argoSdk.ArgoApplication{},
		HistoryId:   historyId,
	}
	// required in order to put into the queue
	aw.Application.Status.OperationState.SyncResult.Revision = "1"
	return &aw
}
