package queue

/*
Queue interface 
*/

type Queue interface {
    Push(job any) error
    Pop() 		  (any, error)
    Peek()        (any, error)
}
