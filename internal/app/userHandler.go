package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrtuuro/blog-aggregator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"Name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

type PostUserHandler struct {
	DB *database.Queries
}

func (h *PostUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	var reqBody parameters
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Insufficiant parameters!")
		return
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
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

type GetUserHandler struct {
	DB *database.Queries
}

func (h *GetUserHandler) Handle(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedId      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(dbPost database.Post) Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedId:      dbPost.FeedID,
	}
}

func databasePostsToPosts(dbPosts []database.Post) []Post {
	var posts []Post
	for _, post := range dbPosts {
		posts = append(posts, databasePostToPost(post))
	}
	return posts
}

type GetPostsForUserHandler struct {
	DB *database.Queries
}

func (h *GetPostsForUserHandler) Handle(w http.ResponseWriter, r *http.Request, user database.User) {
	dbPosts, err := h.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get posts")
		return
	}
	respondWithJSON(w, http.StatusOK, databasePostsToPosts(dbPosts))
}
