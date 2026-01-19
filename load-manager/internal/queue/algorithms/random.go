package algorithms

import (
	"math/rand/v2"
	"sync"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type Random struct {
	jobs  []*queue.Job
	mutex sync.Mutex
}

func (r *Random) Pushs(jobs []*queue.Job) []error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.jobs = append(r.jobs, jobs...)

	errs := make([]error, 0)
	return errs
}

func (r *Random) Pops() ([]*queue.Job, []error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	n := min(r.Len(), MIN_CAPACITY)
	jobs := make([]*queue.Job, n)
	errs := make([]error, n)

	for i := range n {
		idx := rand.IntN(r.Len())
		jobs[i] = r.jobs[idx]
		r.jobs = append(r.jobs[:idx], r.jobs[idx+1:]...)
		errs[i] = nil
	}
	return jobs, errs
}

func (r *Random) Len() int {
	return len(r.jobs)
}
func (r *Random) IsEmpty() bool {
	return r.Len() == 0
}

func NewRand() queue.Queue {
	return &Random{
		jobs: make([]*queue.Job, 0, MIN_CAPACITY),
	}
}
