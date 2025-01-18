package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emiel/chirpy/internal/auth"
)

func loginTask(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input input
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			fmt.Println("error decoding input from body")
		}

		user, err := cfg.db.GetUserByEmail(context.Background(), input.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Incorrect email or password"))
			return
		}

		err = auth.CheckPasswordHash(input.Password, user.HashedPassword)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Incorrect email or password"))
			return
		}

		cleanUser := User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     input.Email}

		encodedUser, err := json.Marshal(cleanUser)
		if err != nil {
			fmt.Println("could not encode user: ", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encodedUser)
	}
}
