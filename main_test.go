package main

import (
	"bytes"
	"net/http"

	"net/http/httptest"
	"testing"
)

func PostRequestTest(t *testing.T) {
	var jsonStr = []byte(`{"amount" : 1654.86,"transaction_time" : "2022-01-10T20:30:00.000Z"}`)

	// reader = strings.NewReader(jsonStr)

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	// rr := httptest.NewRecorder()
	// a.Router.ServeHTTP(rr, req)

	rr := httptest.NewRecorder()
	// res, err := http.DefaultClient.Do(req)
	handler := http.HandlerFunc(postTransaction)
	handler.ServeHTTP(rr, req)
	if err != nil {
		t.Error(err)
	}
	// fmt.Println(res)
	// if status := rr.Code; status != http.StatusCreated {
	// 	t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	// }
}
