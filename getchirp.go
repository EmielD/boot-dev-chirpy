package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func GetChirpTask(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		chirpId := r.PathValue("chirpID")
		fmt.Printf("Extracted chirpID: %s\n", chirpId)

		if len(chirpId) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		parsedChirpId, err := uuid.Parse(chirpId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		chirp, err := cfg.db.GetChirp(r.Context(), parsedChirpId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		type chirpStruct struct {
			Body string `json:"body"`
		}

		encodedChirp, err := json.Marshal(chirpStruct{Body: chirp.Body})
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedChirp)
	}
}
