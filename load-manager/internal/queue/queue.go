package queue


/*
Queue interface 
*/

type Queue interface {
    Pushs([]*Job)   []error
    Pops() 		    ([]*Job, []error)
    Len()           int
    IsEmpty()       bool 
}
