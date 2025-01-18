package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

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

		cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: response.Body, UserID: parsedUuid})

		respondWithJSON(w, 201, response)
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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	result, err := json.Marshal(&payload)
	if err != nil {
		log.Panicf("%v", err)
	}

	w.WriteHeader(code)
	w.Write(result)
}
