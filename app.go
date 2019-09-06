package main

import (
	"encoding/json"
	"net/http"
)

type MenuItem struct {
	Coffee string `json:"name"`
	Type   string `json:"type"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func main() {
	http.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		item := MenuItem{
			Coffee: "Kick Ass",
			Type:   "Iced Coffee",
		}
		w.Header().Add("Content-Type", "application/json")

		json.NewEncoder(w).Encode([]MenuItem{item, item})
	})

	http.ListenAndServe(":5000", nil)
}
