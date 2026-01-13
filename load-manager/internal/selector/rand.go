package selector

import (
	"github.com/sudo-JP/Load-Manager/load-manager/internal/registry"
	"sync"
	"math/rand/v2"
)

type Rand struct {
	mutex sync.Mutex
}

func (r *Rand) SelectNode(nodes []*registry.BackendNode) *registry.BackendNode {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if len(nodes) == 0 {
		return nil 
	}

	node := nodes[rand.IntN(len(nodes))]
	
	return node
}

func NewRand () Selector {
	return &Rand{}
}
