package front_page

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var raw []byte

type MenuItem struct {
	Coffee      string   `json:"coffee"`
	Types       []string `json:"types"`
	Description string   `json:"description"`
}

func Setup(pathPrefix string, router *mux.Router) error {
	subRouter := router.PathPrefix(pathPrefix).Subrouter()

	file, err := os.Open("api/front_page/menu.json")
	if err != nil {
		return err
	}
	defer file.Close()

	raw, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	subRouter.HandleFunc("/menu", MenuHandler)
	return nil
}

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	data := []MenuItem{}
	_ = json.Unmarshal(raw, &data)

	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err)
	}
}
