package algorithms

import (
	"sync"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)


type Stack struct {
	jobs 		[]*queue.Job
	tail 		int 
	mutex 		sync.Mutex
}



func (s *Stack) Pushs(jobs []*queue.Job) []error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	errs := make([]error, 0)

	for _, job := range(jobs) {
		s.jobs = append(s.jobs, job)
	}
	return errs  
}


func (s *Stack) Pops() ([]*queue.Job, []error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	n := min(MIN_CAPACITY, s.Len())
	jobs := make([]*queue.Job, n)
	errs := make([]error, 0)
	for i := range(n) {
		job := s.jobs[s.Len()]
		s.jobs = s.jobs[:s.Len() - 1]
		jobs[i] = job
	}
	return jobs, errs
}


func (s *Stack) Len() int {
	return len(s.jobs)
}
func (s *Stack) IsEmpty() bool {
	return s.Len() == 0
}

func NewStackQueue() queue.Queue {
	return &FCFS{
		jobs: make([]*queue.Job, MIN_CAPACITY),
		tail: 0, 
		capacity: MIN_CAPACITY, 
	}
}

