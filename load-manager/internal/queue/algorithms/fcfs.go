package algorithms

import (
	"errors"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type FCFS struct {
	jobs []any
}

func (q *FCFS) Push(job any) {
	q.jobs = append(q.jobs, job)
}

func (q *FCFS) Pop() (any, error) {
	if len(q.jobs) == 0 {
		return nil, errors.New("queue empty")
	}		

	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job, nil
}

func (q *FCFS) Peek() (any, error) {
	if len(q.jobs) == 0 {
		return nil, errors.New("queue empty")
	}
	return q.jobs[0], nil
}

func NewFCFSQueue() queue.Queue {
	return &FCFS{
		jobs: []any{},
	}
}

