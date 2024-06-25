package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mrtuuro/blog-aggregator/internal/database"
)

type FeedFollow struct {
	ID        uuid.UUID `json:"_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowtoFeedFollow(dbFeedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID:    dbFeedFollow.UserID,
		FeedID:    dbFeedFollow.ID,
	}
}

type FollowFeedHandler struct {
	DB *database.Queries
}

func (h *FollowFeedHandler) Handle(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := h.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create feed follow: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseFeedFollowtoFeedFollow(feedFollow))

}

type GetFollowFeedHandler struct {
	DB *database.Queries
}

func (h *GetFollowFeedHandler) Handle(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := h.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not retrieve feed follows: %v", err))
		return
	}

	feedResp := []FeedFollow{}
	for _, dbFeed := range feedFollows {
		feedResp = append(feedResp, databaseFeedFollowtoFeedFollow(dbFeed))
	}
	respondWithJSON(w, http.StatusOK, feedResp)
}

type DeleteFollowHandler struct {
	DB *database.Queries
}

func (h *DeleteFollowHandler) Handle(w http.ResponseWriter, r *http.Request, user database.User) {

}

type DeleteFeedFromFollowHandler struct {
	DB *database.Queries
}

func (h *DeleteFeedFromFollowHandler) Handle(w http.ResponseWriter, r *http.Request, user database.User) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[1] != "feed_follows" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	feedIdStr := parts[2]
	feedId, err := uuid.Parse(feedIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse feed id: %v", err))
		return
	}
	err = h.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedId,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not delete feed: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})

}
