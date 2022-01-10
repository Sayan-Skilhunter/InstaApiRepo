package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type GetResponse struct {
	Sum   float64
	Avg   float64
	Max   float64
	Min   float64
	Count int
}

type Transaction struct {
	Amount          float64 `json:"amount"`
	TransactionTime string  `json:"transaction_time"`
}

var transaction_list []Transaction

func statistics(w http.ResponseWriter, r *http.Request) {
	fmt.Println(transaction_list)
	var count int
	var sum, max, min float64
	w.Header().Set("content-type", "application/json")
	currentTime := time.Now()

	for _, v := range transaction_list {
		t, err := time.Parse(time.RFC3339Nano, v.TransactionTime)

		if err == nil {
			if t.After(currentTime.Add(time.Second * -60)) {
				sum += v.Amount
				count++
				if v.Amount > max {
					max = v.Amount
				} else if v.Amount < min {
					min = v.Amount
				}
			}
		}
	}
	res := &GetResponse{
		Sum:   sum,
		Avg:   sum / float64(count),
		Max:   max,
		Min:   min,
		Count: count,
	}
	des, _ := json.Marshal(res)
	w.Write(des)
	w.WriteHeader(http.StatusOK)
}

func postTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	w.Header().Set("content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		if strings.Contains(err.Error(), "invalid character") {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	t, err := time.Parse(time.RFC3339Nano, transaction.TransactionTime)
	fmt.Println(t)
	parseError := err
	if parseError == nil {
		fmt.Println("Both timestamps : ", t, time.Now().UTC().Add(time.Second*-60))
		if t.Before(time.Now().UTC().Add(time.Second * -60)) {
			w.Header().Set("content-type", "application/json")
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)

			return
		} else if t.After(time.Now().UTC()) {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
	}
		
		transaction_list = append(transaction_list, transaction)
		fmt.Println(transaction_list)
		w.WriteHeader(http.StatusCreated)

}

func deleteTransactions(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err == io.EOF {
		transaction_list = transaction_list[:0]
		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/statistics", statistics).Methods("GET")
	router.HandleFunc("/transactions", postTransaction).Methods("POST")
	router.HandleFunc("/transactions", deleteTransactions).Methods("DELETE")

	fmt.Println("Server is starting on port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}

}
