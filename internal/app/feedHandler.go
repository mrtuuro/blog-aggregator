package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrtuuro/blog-aggregator/internal/database"
)

type Feed struct {
	ID        uuid.UUID `json:"_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedtoFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		URL:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}

type PostFeedHandler struct {
	DB *database.Queries
}

func (h *PostFeedHandler) Handle(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error decoding JSON: %v", err))
        return
	}

	feed, err := h.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Cannot create feed: %v", err))
        return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedtoFeed(feed))

}


type GetFeedsHandler struct {
    DB *database.Queries
}

func (h *GetFeedsHandler) Handle(w http.ResponseWriter, r *http.Request, user database.User) {
    
    dbFeeds, err := h.DB.GetFeeds(r.Context())
    if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Cannot retrieve feeds: %v", err))
        return
    }

    feeds := []Feed{}
    for _, feed := range dbFeeds {
        returnFeed := databaseFeedtoFeed(feed)
        feeds = append(feeds, returnFeed)
    }

	respondWithJSON(w, http.StatusCreated, feeds)

}
