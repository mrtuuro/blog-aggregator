package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrtuuro/blog-aggregator/internal/database"
)

type PostUserHandler struct {
	DB *database.Queries
}

func (h *PostUserHandler) Handle(w http.ResponseWriter, r *http.Request) {

	var createUser database.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
		respondWithError(w, http.StatusBadRequest, "Insufficiant parameters!")
	}
	now := time.Now().UTC()
	createUser.ID = uuid.New()
	createUser.CreatedAt = now
	createUser.UpdatedAt = now
	user, err := h.DB.CreateUser(r.Context(), createUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user onto DB!")
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"user": user})
}
