package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/emiel/chirpy/internal/auth"
	"github.com/emiel/chirpy/internal/database"
	"github.com/google/uuid"
)

type input struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func createUserTask(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input input
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			fmt.Println("error decoding input from body")
		}

		hashedPassword, err := auth.HashPassword(input.Password)
		if err != nil {
			fmt.Println("error hashing password")
		}

		dbUser, err := cfg.db.CreateUser(r.Context(),
			database.CreateUserParams{
				Email:          input.Email,
				HashedPassword: hashedPassword})
		if err != nil {
			fmt.Println("could not create user")
		}

		user := User{
			Id:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		}

		json, err := json.Marshal(user)
		if err != nil {
			fmt.Println("could not marshal user info")
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		w.Write(json)

	}
}
