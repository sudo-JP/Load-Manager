package worker

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sudo-JP/Load-Manager/load-manager/internal/grpc"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/queue"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/registry"
	"github.com/sudo-JP/Load-Manager/load-manager/internal/selector"


	"log"
)

type LoadBalancingStrategy int 
const (
	Mixed 					LoadBalancingStrategy = iota
	PerResource
	PerOperation
	PerResourceAndOperation
)

type Worker struct {
	queue 		queue.Queue
	registry 	*registry.Registry
	selector 	selector.Selector	
	clients 	map[string]*grpc.BackendClient // key is host:port
	clientsMut 	sync.RWMutex	
	stopCh 		chan struct{}
	workers 	int 
	strategy 	LoadBalancingStrategy
}


func groupByResource(jobs []*queue.Job) map[queue.JobType][]*queue.Job {
	grouped := make(map[queue.JobType][]*queue.Job)

	for _, job := range jobs {
		// error case 
		if job == nil {
			continue
		}
			
		grouped[job.Resource] = append(grouped[job.Resource], job)
	}
	return grouped
}

func groupByCRUD(jobs []*queue.Job) map[queue.Operation][]*queue.Job {
	grouped := make(map[queue.Operation][]*queue.Job)

	for _, job := range jobs {
		if job == nil {
			continue
		}

		grouped[job.CRUD] = append(grouped[job.CRUD], job)
}
	return grouped
}


func (w *Worker) sendToBackend(node *registry.BackendNode, resource queue.JobType, 
	crud queue.Operation, jobs []*queue.Job) {
	if len(jobs) == 0 {
		return 
	}

	log.Printf("Sending %d %v %v jobs to node %s:%d",
		len(jobs), resource, crud, node.Host, node.Port)
	
	switch resource {
	case queue.User: 
		w.sendUserJobs(node, crud, jobs)
	case queue.Product: 
		w.sendProductJobs(node, crud, jobs)
	case queue.Order: 
		w.sendOrdersJobs(node, crud, jobs)
	}
}

func (w *Worker) mixedStat(jobs []*queue.Job) error {
	node := w.selector.SelectNode(w.registry.All())
	if node == nil {
		return errors.New("no available nodes")
	}
	// just optimization for job same type and operation to be tgt  
	grouped := groupByResource(jobs)		
	for resource, resourceJob := range grouped {
		groupedCRUD := groupByCRUD(resourceJob)
		for crud, crudJobs := range groupedCRUD {
			w.sendToBackend(node, resource, crud, crudJobs)
		}
	}
	return nil 
}

func (w *Worker) perOperationStrat(jobs []*queue.Job) {
	groupedCRUD := groupByCRUD(jobs)
	for crud, crudJobs := range groupedCRUD {
		node := w.selector.SelectNode(w.registry.All())

		if node == nil {
			continue
		}

		grouped := groupByResource(crudJobs)
		for resource, resourceJobs := range grouped {
			w.sendToBackend(node, resource, crud, resourceJobs)
		}
	}	
}

func (w *Worker) perResourceStrat(jobs []*queue.Job) {
	groupedResource := groupByResource(jobs)
	// optimization 
	for resource, resourceJobs := range groupedResource {
		node := w.selector.SelectNode(w.registry.All())
		if node == nil {
			continue 
		}			
		
		grouped := groupByCRUD(resourceJobs)
		for crud, crudJobs := range grouped {
			w.sendToBackend(node, resource, crud, crudJobs)
		}
	}	
}

func (w *Worker) perOperationAndResouceStrat(jobs []*queue.Job) error {
	groupedCRUD := groupByCRUD(jobs)	

	// per crud 
	for crud, crudJobs := range groupedCRUD {

		// per resource 
		grouped := groupByResource(crudJobs)
		for resource, resourceJobs := range grouped {

			// each type, we get a new node and send to backend 
			node := w.selector.SelectNode(w.registry.All())

			if node == nil {
				return errors.New("no available nodes")
			}

			w.sendToBackend(node, resource, crud, resourceJobs)

		}

	}
	return nil 
}

func (w *Worker) run() {
	for {
		select {
		case <- w.stopCh: 
			return 
		default:
		}

		jobs, errs := w.queue.Pops() 
		if len(errs) > 0 {
			log.Printf("Errors popping from queue %v", errs)
		}

		if len(jobs) == 0 {
			time.Sleep(10 * time.Millisecond) // sleep so no cpu waste 
			continue
		}

		// pick strategy
		switch w.strategy {

		// One node for all operations and resource
		case Mixed: 
			err := w.mixedStat(jobs)	
			if err != nil {
				log.Println(err)
			}

		// One node per operation 
		case PerOperation: 
			w.perOperationStrat(jobs)

		// One node per resource 
		case PerResource:
			w.perResourceStrat(jobs)

		// One node per resouce and operation 
		case PerResourceAndOperation: 
			err := w.perOperationAndResouceStrat(jobs)
			if err != nil {
				log.Println(err)	
			}	
		}
	}
}

func (w *Worker) getClient(node *registry.BackendNode) (*grpc.BackendClient, error) {
	addr := fmt.Sprintf("%s:%d", node.Host, node.Port)

	w.clientsMut.RLock()
	client, ok := w.clients[addr]
	w.clientsMut.RUnlock()

	if ok {
		return client, nil 
	}

	// New client connection
	w.clientsMut.Lock()
	defer w.clientsMut.Unlock()

	client, err := grpc.NewBackendClient(addr)
	if err != nil {
		return nil, err
	}

	w.clients[addr] = client
	return client, nil
}

func (w *Worker) sendUserJobs(node *registry.BackendNode, 
	crud queue.Operation, jobs []*queue.Job) {
	switch crud {
	case queue.Create: 
		w.CreateUsers(node, jobs)
	case queue.Read:

	}
}

func (w *Worker) sendProductJobs(node *registry.BackendNode, 
	crud queue.Operation, jobs []*queue.Job) {
	// TODO: call gRPC
}

func (w *Worker) sendOrdersJobs(node *registry.BackendNode, 
	crud queue.Operation, jobs []*queue.Job) {
	// TODO: call gRPC
}

func (w *Worker) Stop() {

	close(w.stopCh)
}

func NewWorker(q queue.Queue, reg *registry.Registry, selector selector.Selector, 
	clients map[string]*grpc.BackendClient, workers int, strat LoadBalancingStrategy) *Worker {
	w := &Worker{
		queue: 		q, 
		registry: 	reg, 
		clients: 	clients, 
		selector: 	selector,
		workers: 	workers, 
		stopCh: 	make(chan struct{}), 
		strategy: 	strat, 
	}

	for range workers {
		go w.run()
	}
	return w
}
