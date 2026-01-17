package algorithms

import (
	"testing"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

func TestFCFS_PushPop(t *testing.T) {
	q := NewFCFSQueue()
	testSize := 100

	arr := make([]*queue.Job, testSize)
	for i := range(testSize) {
		job := &queue.Job{ID: i}
		arr[i] = job
	}

	q.Pushs(arr)
	
	if q.Len() != testSize {
		t.Errorf("Expected %d jobs, got %d", testSize, q.Len())
	}

	popped := 0 

	for q.Len() != 0 {
		jobs, errs := q.Pops()
		
		for _, err := range(errs) {
			if err != nil {
				t.Errorf("Queue not supposed to error %v", err)
			}
		}

		popped += len(jobs)
	}

	if popped != testSize {
		t.Errorf("Did not pop enough elements %d", popped)
	}
}
