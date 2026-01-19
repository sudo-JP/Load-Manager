package algorithms

import (
	"errors"
	"sync"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

// LJF - Longest Job First (max-heap by payload size)
type LJF struct {
	jobs  []*JobPriority
	mutex sync.Mutex
}

func (l *LJF) parent(idx int) int {
	return (idx - 1) >> 1
}

func (l *LJF) left_child(idx int) int {
	return 2*idx + 1
}

func (l *LJF) right_child(idx int) int {
	return 2*idx + 2
}

// Bubble up last index to maintain max-heap property
func (l *LJF) bubbleUp() {
	i := len(l.jobs) - 1
	for i > 0 {
		parent := l.jobs[l.parent(i)]
		child := l.jobs[i]

		// Max-heap: parent should be >= child
		if parent.priority < child.priority {
			l.jobs[i] = parent
			l.jobs[l.parent(i)] = child
			i = l.parent(i)
		} else {
			break
		}
	}
}

func (l *LJF) push(job *queue.Job) error {
	if job == nil {
		return errors.New("nil job")
	}

	node := &JobPriority{
		priority: len(job.Payload),
		job:      job,
	}
	l.jobs = append(l.jobs, node)
	l.bubbleUp()

	return nil
}

func (l *LJF) Pushs(jobs []*queue.Job) []error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	errs := make([]error, len(jobs))
	for i, job := range jobs {
		errs[i] = l.push(job)
	}

	return errs
}

func (l *LJF) bubbleDown(idx int) {
	size := len(l.jobs)
	for {
		largest := idx
		left := l.left_child(idx)
		right := l.right_child(idx)

		// Check left child - max-heap: find largest
		if left < size && l.jobs[left].priority > l.jobs[largest].priority {
			largest = left
		}

		// Check right child
		if right < size && l.jobs[right].priority > l.jobs[largest].priority {
			largest = right
		}

		// If current is largest, heap property satisfied
		if largest == idx {
			break
		}

		// Swap and continue
		l.jobs[idx], l.jobs[largest] = l.jobs[largest], l.jobs[idx]
		idx = largest
	}
}

func (l *LJF) pop() (*queue.Job, error) {
	if l.Len() == 0 {
		return nil, errors.New("empty queue")
	}

	// Get root (largest)
	job := l.jobs[0].job

	// Move last to root
	lastIdx := len(l.jobs) - 1
	l.jobs[0] = l.jobs[lastIdx]
	l.jobs = l.jobs[:lastIdx]

	// Restore heap property
	if len(l.jobs) > 0 {
		l.bubbleDown(0)
	}

	return job, nil
}

func (l *LJF) Pops() ([]*queue.Job, []error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	n := min(MIN_CAPACITY, l.Len())
	jobs := make([]*queue.Job, n)
	errs := make([]error, n)

	for i := range n {
		jobs[i], errs[i] = l.pop()
	}

	return jobs, errs
}

func (l *LJF) Len() int {
	return len(l.jobs)
}

func (l *LJF) IsEmpty() bool {
	return len(l.jobs) == 0
}

func NewLJF() queue.Queue {
	return &LJF{
		jobs: make([]*JobPriority, 0, MIN_CAPACITY),
	}
}
