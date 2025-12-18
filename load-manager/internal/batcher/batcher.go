package batcher

import (
	"sync"
	"time"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type Batcher struct {
	queue    	queue.Queue
	users    	[]*queue.Job
	products 	[]*queue.Job
	orders   	[]*queue.Job

	batchSize 	int
	timeout   	time.Duration

	mutex  		sync.Mutex
	timer  		*time.Timer
	stopCh 		chan struct{}
}

func (b *Batcher) AddUser(job *queue.Job) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.users = append(b.users, job)
	if len(b.users) >= b.batchSize {
		b.flushUsersLocked()
	}
}

func (b *Batcher) AddProduct(job *queue.Job) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.products = append(b.products, job)
	if len(b.products) >= b.batchSize {
		b.flushUsersLocked()
	}
}

func (b *Batcher) AddOrder(job *queue.Job) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.orders = append(b.orders, job)
	if len(b.orders) >= b.batchSize {
		b.flushUsersLocked()
	}
}

func (b *Batcher) flush() {
	b.mutex.Lock()

	users := b.users
	products := b.products
	orders := b.orders

	b.users = make([]*queue.Job, 0, b.batchSize)
	b.products = make([]*queue.Job, 0, b.batchSize)
	b.orders = make([]*queue.Job, 0, b.batchSize)

	b.mutex.Unlock()

	b.groupAndPush(users)
	b.groupAndPush(products)
	b.groupAndPush(orders)
}

func (b *Batcher) flushUsersLocked() {
	users := b.users 
	b.users = make([]*queue.Job, 0, b.batchSize)
	b.mutex.Unlock()
	b.groupAndPush(users)
	b.mutex.Lock()
}

func (b *Batcher) groupAndPush(jobs []*queue.Job) {
	if len(jobs) == 0 {
		return
	}

	creates := make([]*queue.Job, 0, b.batchSize)
	reads := make([]*queue.Job, 0, b.batchSize)
	updates := make([]*queue.Job, 0, b.batchSize)
	deletes := make([]*queue.Job, 0, b.batchSize)

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
		b.queue.Pushs(creates)
	}
	if len(reads) > 0 {
		b.queue.Pushs(reads)
	}
	if len(updates) > 0 {
		b.queue.Pushs(updates)
	}
	if len(deletes) > 0 {
		b.queue.Pushs(deletes)
	}
}

func (b *Batcher) run() {
	for {
		select {
		case <-b.timer.C:
			b.flush()
			b.timer.Reset(b.timeout)
		case <-b.stopCh:
			return
		}
	}
}

func NewBatcher(queue queue.Queue, batchSize int, timeout time.Duration) *Batcher {
	b := &Batcher{
		queue:     queue,
		batchSize: batchSize,
		timeout:   timeout,
		users:     make([]*queue.Job, 0, batchSize),
		products:  make([]*queue.Job, 0, batchSize),
		orders:    make([]*queue.Job, 0, batchSize),
		timer:     time.NewTimer(timeout),
		stopCh:    make(chan struct{}),
	}

	go b.run()
	return b
}
