package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrtuuro/blog-aggregator/internal/auth"
	"github.com/mrtuuro/blog-aggregator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"Name"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
	}
}

type PostUserHandler struct {
	DB *database.Queries
}

func (h *PostUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"Name"`
	}

	var reqBody parameters
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Insufficiant parameters!")
	}
	if reqBody.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name can not be empty!")
		return
	}

	user, err := h.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqBody.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user onto DB!")
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

type GetUserHandler struct {
	DB *database.Queries
}

func (h *GetUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	dbUser, err := h.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving user!")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseUserToUser(dbUser))
}
