package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func GetChirpsTask(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		chirps, err := cfg.db.GetChirps(context.Background())
		if err != nil {
			log.Panicf("%s", err.Error())
		}

		type JsonChirp struct {
			Id        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserId    uuid.UUID `json:"user_id"`
		}

		var jsonChirps []JsonChirp
		for _, chirp := range chirps {
			jsonChirps = append(jsonChirps, JsonChirp{
				Id:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserId:    chirp.UserID,
			})
		}

		encodedChirps, err := json.Marshal(jsonChirps)
		if err != nil {
			log.Panicf("%s", err.Error())
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		w.Write(encodedChirps)
	}
}
