package queue

import "time"

var idCounter int = -1

type JobType int
type Operation int 

const (
	User 	JobType = iota
	Product 
	Order
)

const (
	Create Operation = iota	
	Read 
	Update 
	Delete 
)


type Job struct {
	ID 			int
	Resource 	JobType	
	CRUD 		Operation
	Payload 	[]byte
	Priority 	int 
	CreatedAt 	time.Time
}

func GetID() int {
	idCounter++
	return idCounter
}

