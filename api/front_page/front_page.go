package front_page

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

var file []byte

type MenuItem struct {
	Coffee string `json:"coffee"`
	Type   string `json:"type"`
}

func Setup(pathPrefix string, router *mux.Router) error {
	var err error
	subRouter := router.PathPrefix(pathPrefix).Subrouter()

	file, err = ioutil.ReadFile("menu.json")
	if err != nil {
		return err
	}

	subRouter.HandleFunc("/menu", MenuHandler)
	return nil
}

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	data := []MenuItem{}
	_ = json.Unmarshal([]byte(file), &data)

	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err)
	}
}
