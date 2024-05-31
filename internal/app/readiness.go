package app

import (
	"net/http"
)

type ReadinessHandler struct {
}

func (h *ReadinessHandler) Handle(w http.ResponseWriter, r *http.Request) {
    respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


type ErrorHandler struct {
} 

func (h *ErrorHandler) Handle(w http.ResponseWriter, r *http.Request) {
    respondWithError(w, http.StatusInternalServerError, "Internal Server Error!")

}
