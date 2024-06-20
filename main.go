package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var LOCAL_COUNTER = 0

func pullHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, fmt.Sprintf("%d", LOCAL_COUNTER))
}

func main() {
	go value()
	http.HandleFunc("/pull", pullHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func requestDataFromOtherNodes() string {
	resp, err := http.Get("http://localhost:8080/pull")
	if err != nil {
		fmt.Printf("Error making GET request: %s\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return ""
	}

	return string(body)
}

func merge() {
	dataFromOtherNode := requestDataFromOtherNodes()
	actualDataInInt, err := strconv.ParseInt(dataFromOtherNode, 2, 64)
	if err != nil {
		fmt.Errorf("error %w", err)
		return
	}
	if actualDataInInt > int64(LOCAL_COUNTER) {
		LOCAL_COUNTER = int(actualDataInInt)
	}
}

func value() {
	tick := time.Tick(2 * time.Second)
	localOffset := rand.Intn(100)
	for {
		select {
		case <-tick:
			LOCAL_COUNTER = LOCAL_COUNTER + localOffset
			merge()
		}
	}
}
