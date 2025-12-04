package queue

// Priority at element 0 
// drops down at n
type MLFQ struct {
	queues []Queue
}
