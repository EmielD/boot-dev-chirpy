package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/emiel/chirpy/internal/database"
	"github.com/google/uuid"
)

type response struct {
	Body    string `json:"body"`
	User_Id string `json:"user_id"`
}

type errorResponse struct {
	Error   string `json:"error"`
	User_Id string `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func ChirpsTask(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var response response

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&response)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error(), response.User_Id)
			return
		}

		if len(response.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long", response.User_Id)
			return
		}

		response.Body = checkForProfanity(response.Body)

		parsedUuid, err := uuid.Parse(response.User_Id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing user UUID", response.User_Id)
		}

		chirpDb, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: response.Body, UserID: parsedUuid})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error inserting chirp into db", response.User_Id)
		}

		chirp := Chirp{
			ID:        chirpDb.ID,
			CreatedAt: chirpDb.CreatedAt,
			UpdatedAt: chirpDb.UpdatedAt,
			Body:      chirpDb.Body,
			UserID:    chirpDb.UserID,
		}

		respondWithJSON(w, chirp)
	}
}

func checkForProfanity(body string) (cleanedBody string) {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	body_splitted := strings.Split(body, " ")

	for i, word := range body_splitted {
		if slices.Contains(bannedWords, strings.ToLower(word)) {
			body_splitted[i] = "****"
		}
	}

	return strings.Join(body_splitted, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string, userId string) {
	response := errorResponse{
		Error:   msg,
		User_Id: userId,
	}
	result, err := json.Marshal(&response)
	if err != nil {
		log.Panicf("%v", err)
	}

	w.WriteHeader(code)
	w.Write(result)
}

func respondWithJSON(w http.ResponseWriter, payload interface{}) {
	result, err := json.Marshal(&payload)
	if err != nil {
		log.Panicf("%v", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}
