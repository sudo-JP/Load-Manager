package main 

import ( 
	"fmt"
	"net/http"
    "github.com/sudo-JP/Load-Manager/load-manager/internal/queue/algorithms"
    "github.com/sudo-JP/Load-Manager/load-manager/internal/routes"
)

func main() {
	q := algorithms.NewFCFSQueue()

	q.Push("job1")
	q.Push("job2")

	next, _ := q.Peek()
    fmt.Println("Next job:", next)

    job, _ := q.Pop()
    fmt.Println("Popped:", job)

    job, _ = q.Pop()
    fmt.Println("Popped:", job)
	routes.RegisterRoutes()
	http.ListenAndServe(":8080", nil)
}

