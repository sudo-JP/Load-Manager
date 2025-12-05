package routes

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/users", UserRoutes)
	http.HandleFunc("/orders", OrderRoutes)
	http.HandleFunc("/products", ProductRoutes)
}
