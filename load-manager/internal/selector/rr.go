package selector


import (
	"github.com/sudo-JP/Load-Manager/load-manager/internal/registry"
	"sync"
)

type RoundRobin struct {
	next int 
	mutex sync.Mutex
}

func (rr *RoundRobin) SelectNode(nodes []*registry.BackendNode) *registry.BackendNode {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()
	
	if len(nodes) == 0 {
		return nil 
	}

	node := nodes[rr.next % len(nodes)]
	rr.next++
	
	return node
}

func NewRR () Selector {
	return &RoundRobin{
		next: 0, 
	}
}
