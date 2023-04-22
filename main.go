package main

import (
	"encoding/json"
	"fmt"
	"github.com/AlmasOrazgaliev/halyk-life-task/models"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var CounterChan chan int

func Handler(w http.ResponseWriter, r *http.Request) {
	var request models.Request
	var response models.Response
	w.Header().Set("Content-Type", "application/json")

	// accepting request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	// sending request
	method := request.Method
	url := request.Url
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	req.Header.Add("Authentication", request.Headers.Authentication)
	client := &http.Client{}
	proxyResponse, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	// building response
	w.WriteHeader(http.StatusOK)
	response.Status = proxyResponse.StatusCode

	var headers []string
	for name, _ := range proxyResponse.Header {
		headers = append(headers, name)
	}
	response.Headers = headers

	b, _ := io.ReadAll(proxyResponse.Body)
	response.Length = len(b)

	response.Id = <-CounterChan

	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	counter := 0
	CounterChan = make(chan int)
	go func() {
		for {
			fmt.Println(counter)
			counter++
			CounterChan <- counter
		}

	}()
	router.HandleFunc("/proxy", Handler)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
