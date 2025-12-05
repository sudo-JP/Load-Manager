package routes

import (
	"encoding/json"
	"net/http"
)

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r.Body)
}

func DeleteProductsHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdateProductsHandler(w http.ResponseWriter, r *http.Request) {

}

func CreateProductsHandler(w http.ResponseWriter, r *http.Request) {

}

func ProductRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetProductsHandler(w, r)
	case http.MethodPost:
		CreateProductsHandler(w, r)
	case http.MethodPut:
		UpdateProductsHandler(w, r)
	case http.MethodDelete:
		DeleteProductsHandler(w, r)
	default:
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
	}
}
