package algorithms

import (
	"errors"
	"sync"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type JobPriority struct {
	priority	int 
	job			*queue.Job
}


type SJF struct {
	jobs 		[]*JobPriority
	mutex 		sync.Mutex
}

func (s *SJF) parent(idx int) int {
	return (idx - 1) >> 1
}

func (s *SJF) left_child(idx int) int {
	return 2 * idx + 1
}

func (s *SJF) right_child(idx int) int {
	return 2 * idx + 2
}

// Always bubble up last index to top ish
func (s *SJF) bubbleUp() {
	i := len(s.jobs) - 1
	for i > 0 {
		parent := s.jobs[s.parent(i)]
		child := s.jobs[i]

		if parent.priority > child.priority {
			s.jobs[i] = parent
			s.jobs[s.parent(i)] = child
		} else {
			break
		}
	}
}

func (s *SJF) push(job *queue.Job) error {
	if job == nil {
		return errors.New("nil job")
	}


	node := &JobPriority{
		priority: len(job.Payload), 
		job: job,
	}
	s.jobs = append(s.jobs, node)

	return nil
}

func (s *SJF) Pushs(jobs []*queue.Job) []error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	errs := make([]error, len(jobs))
	for i, job := range(jobs) {
		errs[i] = s.push(job)	
	}

	return errs
}

func (s *SJF) pop() (*queue.Job, error) {
	if s.Len() == 0 {
		return nil, errors.New("empty queue")
	}

	i := 0 
	job := s.jobs[i].job
	
	for i < s.Len() {
		left := s.left_child(i)
		right := s.right_child(i)

		curr_p := s.jobs[i]
		left_p := s.jobs[left]
		right_p := s.jobs[right]

		if curr_p.priority < left_p.priority && curr_p.priority < right_p.priority {
			break
		} else if left_p.priority < right_p.priority {
			s.jobs[i] = s.jobs[left]
			s.jobs[left] = curr_p
		} else {
			s.jobs[i] = s.jobs[right]
			s.jobs[right] = curr_p
		}
		
	}
	return job, nil
}

func (s *SJF) Pops() ([]*queue.Job, []error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	n := min(MIN_CAPACITY, s.Len())
	jobs := make([]*queue.Job, n)
	errs := make([]error, n)

	for i := range(n) {
		jobs[i], errs[i] = s.pop()
	}

	return jobs, errs
}

func (s *SJF) Len() int {
	return len(s.jobs)
}

func (s *SJF) IsEmpty() bool {
	return len(s.jobs) == 0
}


func NewSJF() queue.Queue {
	return &SJF{
		jobs: make([]*JobPriority, MIN_CAPACITY),
	}
}

