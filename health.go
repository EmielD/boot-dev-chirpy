package main

import (
	"fmt"
	"net/http"
)

func HealthTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	fmt.Println("Got a new request for /healthz")

	w.Write([]byte("OK"))
}
