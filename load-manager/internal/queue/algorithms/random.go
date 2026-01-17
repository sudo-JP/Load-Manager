package algorithms

import (
	"sync"
	//"math/rand/v2"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type Random struct {
	jobs 		[]*queue.Job
	mutex 		sync.Mutex
}

func (r *Random) Pushs(jobs []*queue.Job) []error {
	r.mutex.Lock()	
	defer r.mutex.Unlock()
	r.jobs = append(r.jobs, jobs...)

	errs := make([]error, 0)
	return errs
}

func (r *Random) Pops() ([]*queue.Job, []error) {
	return make([]*queue.Job, 0), make([]error, 0)
}

func (r *Random) Len() int {
	return len(r.jobs)
}
func (r *Random) IsEmpty() bool {
	return r.Len() == 0
}

func NewRand() queue.Queue {
	return &Random {
		jobs: make([]*queue.Job, MIN_CAPACITY),
	}
}
