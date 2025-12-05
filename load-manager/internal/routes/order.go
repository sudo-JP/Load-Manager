package routes

import (
	"encoding/json"
	"net/http"
)

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.Body)
}

func DeleteOrdersHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdateOrdersHandler(w http.ResponseWriter, r *http.Request) {

}

func CreateOrdersHandler(w http.ResponseWriter, r *http.Request) {

}

func OrderRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetOrdersHandler(w, r)
	case http.MethodPost:
		CreateOrdersHandler(w, r)
	case http.MethodPut:
		UpdateOrdersHandler(w, r)
	case http.MethodDelete:
		DeleteOrdersHandler(w, r)
	default:
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
	}
}
