
package stress

import (
    "net/http"
    "sync"
    "time"
)

func sendRequest(route string, numReqs int) {
    client := &http.Client{}
    var wg sync.WaitGroup
    wg.Add(numReqs)

    for i := 0; i < numReqs; i++ {
        go func() {
            defer wg.Done()

            req, _ := http.NewRequest("GET", "http://localhost:8080"+route, nil)
            start := time.Now()
            resp, err := client.Do(req)
            _ = time.Since(start)

            if err != nil {
                return
            }
            resp.Body.Close()
        }()
    }

    wg.Wait()
}

func Stress(numReqs int) {
    routes := []string{"/users", "/products", "/orders"}
    var wg sync.WaitGroup
    wg.Add(len(routes))

    for _, route := range routes {
        go func(r string) {
            defer wg.Done()
            sendRequest(r, numReqs)
        }(route)
    }

    wg.Wait()
}
