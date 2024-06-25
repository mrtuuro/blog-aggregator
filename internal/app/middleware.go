package app

import (
	"net/http"

	"github.com/mrtuuro/blog-aggregator/internal/auth"
	"github.com/mrtuuro/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (a *Application) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		dbUser, err := a.Cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error retrieving user!")
			return
		}
		handler(w, r, dbUser)
	}
}
