package routes

import (
	"encoding/json"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) 
	}
}

func DeleteUsersHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdateUsersHandler(w http.ResponseWriter, r *http.Request) {

}

func CreateUsersHandler(w http.ResponseWriter, r *http.Request) {

}

func UserRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetUsersHandler(w, r)
	case http.MethodPost: 
		CreateUsersHandler(w, r)
	case http.MethodDelete: 
		DeleteUsersHandler(w, r)
	case http.MethodPut: 
		UpdateUsersHandler(w, r)
	default: 
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
	}
}
