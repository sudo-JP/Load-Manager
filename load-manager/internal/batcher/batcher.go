package batcher

import (
	"sync"
	"time"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type Batcher struct {
	Queue    queue.Queue

	Users    []*queue.Job
	Products []*queue.Job
	Orders   []*queue.Job

	BatchSize int
	Timeout   time.Duration

	Mutex  sync.Mutex
	Timer  *time.Timer
	StopCh chan struct{}
}


func (b *Batcher) AddUser(job *queue.Job) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	b.Users = append(b.Users, job)
}

func (b *Batcher) AddProduct(job *queue.Job) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	b.Products = append(b.Products, job)
}

func (b *Batcher) AddOrder(job *queue.Job) {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	b.Orders = append(b.Orders, job)
}

func (b *Batcher) flush() {
	b.Mutex.Lock()

	users := b.Users 
	products := b.Products 
	orders := b.Orders

	b.Users = make([]*queue.Job, 0, b.BatchSize)
	b.Products = make([]*queue.Job, 0, b.BatchSize)
	b.Orders = make([]*queue.Job, 0, b.BatchSize)

	b.Mutex.Unlock()

	b.groupAndPush(users)
	b.groupAndPush(products)
	b.groupAndPush(orders)
}

func (b *Batcher) groupAndPush(jobs []*queue.Job) {
	if len(jobs) == 0 {
		return 
	}

	creates := make([]*queue.Job, b.BatchSize)
	reads := make([]*queue.Job, b.BatchSize)
	updates := make([]*queue.Job, b.BatchSize)
	deletes := make([]*queue.Job, b.BatchSize)

	for _, job := range jobs {
		switch job.CRUD {
		case queue.Create: 
			creates = append(creates, job)
		case queue.Read: 
			reads = append(reads, job)
		case queue.Update: 
			updates = append(updates, job)	
		case queue.Delete: 
			deletes = append(deletes, job)
		}
	}

	if len(creates) > 0 {
		b.Queue.Pushs(creates)
	}
	if len(reads) > 0 {
		b.Queue.Pushs(reads)
	}
	if len(updates) > 0 {
		b.Queue.Pushs(updates)
	}
	if len(deletes) > 0 {
		b.Queue.Pushs(deletes)
	}
}


// If timer runs out, flush everything 
func (b *Batcher) run() {
	for {
		select {
		case <-b.Timer.C:
			b.flush()
			b.Timer.Reset(b.Timeout)
		case <-b.StopCh:
			return
		}
	}
}

func NewBatcher(queue queue.Queue, batchSize int, timeout time.Duration) *Batcher {
	b := &Batcher{
		Queue:     queue,
		BatchSize: batchSize,
		Timeout:   timeout,
		Users:     make([]*queue.Job, 0, batchSize),
		Products:  make([]*queue.Job, 0, batchSize),
		Orders:    make([]*queue.Job, 0, batchSize),
		Timer:     time.NewTimer(timeout),
		StopCh:    make(chan struct{}),
	}

	go b.run()
	return b
}
