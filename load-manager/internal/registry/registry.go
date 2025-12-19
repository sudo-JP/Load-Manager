package registry

import (
	"sync"
	"time"
)

type BackendNode struct {
	ID 				int 
	Host 			string 
	Port 			int	
	Health 			bool
	ActiveReqCount 	int32	
}

type Registry struct {
	Nodes 	[]*BackendNode
	mutex 	sync.RWMutex
	nextID 	int // For setting backend id 
}

func (r *Registry) Add(host string, port int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	node := BackendNode{
		ID: 			r.nextID,
		Host: 			host, 
		Health: 		false, 
		Port: 			port,
		ActiveReqCount: 0,
	}
	r.nextID++
	r.Nodes = append(r.Nodes, &node)
}

func (r *Registry) All() []*BackendNode {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]*BackendNode, len(r.Nodes))
	copy(result, r.Nodes)
	return result
}

func (r *Registry) SetHealth(id int, healthy bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i := range(r.Nodes) {
		if r.Nodes[i].ID == id {
			r.Nodes[i].Health = healthy
			return 
		}
	}
}

func (r *Registry) Remove(id int) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, node := range(r.Nodes) {
		if node.ID == id {
			// Slicing to remove 
			r.Nodes = append(r.Nodes[:i], r.Nodes[i+1:]...)
			return true 
		}
	}
	return false 
}

func (r *Registry) HealthCheckLoop() {
	for {
		time.Sleep(10 * time.Second) 

		nodes := r.All()
		for _, node := range nodes {
			healthy := r.checkHealth(node)
			r.SetHealth(node.ID, healthy)
		}
	}
}

func (r *Registry) checkHealth(node *BackendNode) bool {
	return false 
}
