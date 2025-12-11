package routes

import (
	"net/http"
	"encoding/json"
)

func RegisterRoutes() {
	http.HandleFunc("/data", GetDataHandler)
}

func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) 
	}
}
