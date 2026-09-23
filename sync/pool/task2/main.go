package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type RequestData struct {
	Name  string            `json:"name"`
	Items map[string]string `json:"items"`
}

var dataPool = sync.Pool{
	New: func() any {
		return &RequestData{
			Items: make(map[string]string),
		}

	},
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	data := dataPool.Get().(*RequestData)
	defer func() {
		data.Reset()
		dataPool.Put(data)
	}()
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(data)

}

func (r *RequestData) Reset() {
	r.Name = ""
	clear(r.Items)
}
func main() {

	http.HandleFunc("/", handleRequest)
	fmt.Println("Server started at :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
