package batcher

import (
	"sync"
	"time"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
)

type PendingItem struct {
	Data 		any
	Timestamp 	time.Time
}

type Batcher struct {
	queue *queue.Queue
	
	users 			[]*PendingItem
	products 		[]*PendingItem
	orders 			[]*PendingItem

	batchSize 		int 
	timeout 		time.Duration

	mutex 			sync.Mutex

	userTimer 		*time.Timer
	productTimer 	*time.Timer
	orderTimer 		*time.Timer
}

func (b *Batcher) AddUser(data any) {
	
}


func (b *Batcher) AddProduct(data any) {
	
}

func (b *Batcher) AddOrder(data any) {
	
}
