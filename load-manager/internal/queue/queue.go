package queue

/*
Queue interface 
*/

type Queue interface {
    Push(job any) 
    Pop() 		  (any, error)
    Peek()        (any, error)
}
