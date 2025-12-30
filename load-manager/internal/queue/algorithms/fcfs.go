package algorithms

import (
	"errors"
	"sync"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

const MIN_CAPACITY = 128 

type FCFS struct {
	jobs 		[]*queue.Job
	head 		int 
	tail 		int 
	capacity 	int 
	size 		int 
	mutex 		sync.Mutex
}

func (q *FCFS) resizeQueue(oldCap int, tempArr []*queue.Job) {
	idx := 0 
	for i := q.head; i != q.tail; i = (i + 1) % oldCap {
		tempArr[idx] = q.jobs[i]
		idx++
	}
	q.head = 0 
	q.tail = idx
	q.jobs = tempArr
}

func (q *FCFS) doubleQueue() {
	// Double array 
	oldCap := q.capacity
	q.capacity <<= 1
	tempArr := make([]*queue.Job, q.capacity)
	q.resizeQueue(oldCap, tempArr)
}

func (q *FCFS) Pushs(jobs []*queue.Job) []error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	errs := make([]error, len(jobs))
	for i, job := range(jobs) {
		errs[i] = q.push(job)	
	}
	return errs  
}

func (q *FCFS) push(job *queue.Job) error {
	// Calculate length 
	if job == nil {
		return errors.New("nil job")
	}
	if q.size == q.capacity {
		q.doubleQueue()
	}
	q.size++ 
	q.jobs[q.tail] = job
	q.tail = (q.tail + 1) % q.capacity

	return nil 
}

func (q *FCFS) halfQueue() {
	oldCap := q.capacity
	q.capacity >>= 1
	tempArr := make([]*queue.Job, q.capacity)
	q.resizeQueue(oldCap, tempArr)
}

func (q *FCFS) Pops() ([]*queue.Job, []error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	n := min(MIN_CAPACITY, q.Len())
	jobs := make([]*queue.Job, n)
	errs := make([]error, n)
	for i := range(n) {
		jobs[i], errs[i] = q.pop()
	}
	return jobs, errs
}

func (q *FCFS) pop() (*queue.Job, error) {
	if q.Len() == 0 {
		return nil, errors.New("pop empty on FCFS queue")
	}
	job := q.jobs[q.head]
	q.jobs[q.head] = nil 
	q.head = (q.head + 1) % q.capacity
	q.size--
	if q.size <= q.capacity >> 2 && q.capacity > MIN_CAPACITY {
		q.halfQueue()
	}

	return job, nil
}

func (q *FCFS) Len() int {
	return q.size 
}
func (q *FCFS) IsEmpty() bool {
	return q.size == 0
}

func NewFCFSQueue() queue.Queue {
	return &FCFS{
		jobs: make([]*queue.Job, MIN_CAPACITY),
		head: 0,
		tail: 0, 
		size: 0, 
		capacity: MIN_CAPACITY, 
	}
}

