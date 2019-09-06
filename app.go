package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type MenuItem struct {
	Coffee string `json:"coffee"`
	Type   string `json:"type"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func main() {
	file, _ := ioutil.ReadFile("menu.json")

	data := []MenuItem{}
	_ = json.Unmarshal([]byte(file), &data)
	fmt.Println(data)

	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("running in debug mode")
		port = "5000"
	}

	http.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)

		w.Header().Add("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Println(err)
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
