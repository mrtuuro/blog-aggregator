package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrtuuro/blog-aggregator/internal/config"
)

type Application struct {
	HTTPServer *http.Server
	Cfg        *config.Config
	Context    context.Context
}

func New(cfg *config.Config) *Application {

	application := &Application{
		Cfg:     cfg,
		Context: context.Background(),
	}

	application.HTTPServer = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: application.NewRouter(),
	}

	return application
}

func (a *Application) Start() error {
	errs := make(chan error, 1)
	go func() {
		errs <- a.HTTPServer.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	return <-errs
}

func (a *Application) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := &ReadinessHandler{}
		handler.Handle(w, r)
	}))

	mux.Handle("GET /err", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := &ErrorHandler{}
		handler.Handle(w, r)
	}))

	userHandler := &GetUserHandler{
		DB: a.Cfg.DB,
	}
	mux.Handle("GET /user", a.middlewareAuth(userHandler.Handle))

	mux.Handle("POST /users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		postUserHandler := &PostUserHandler{
			DB: a.Cfg.DB,
		}
		postUserHandler.Handle(w, r)
	}))

	postFeedHandler := &PostFeedHandler{
		DB: a.Cfg.DB,
	}
	mux.Handle("POST /feed", a.middlewareAuth(postFeedHandler.Handle))

	getFeedsHandler := &GetFeedsHandler{
		DB: a.Cfg.DB,
	}
	mux.Handle("GET /feeds", a.middlewareAuth(getFeedsHandler.Handle))

	feedFollowHandler := &FollowFeedHandler{
		DB: a.Cfg.DB,
	}
	mux.Handle("POST /feed_follows", a.middlewareAuth(feedFollowHandler.Handle))

	getFeedFollowHandler := &GetFollowFeedHandler{
		DB: a.Cfg.DB,
	}
	mux.Handle("GET /feed_follows", a.middlewareAuth(getFeedFollowHandler.Handle))

    deleteFeedFollowHandler := &DeleteFeedFromFollowHandler{
        DB: a.Cfg.DB,
    }
    mux.Handle("DELETE /feed_follows/{id}", a.middlewareAuth(deleteFeedFollowHandler.Handle))

    getPostsForUserHandler := &GetPostsForUserHandler{
        DB: a.Cfg.DB,
    }
    mux.Handle("GET /posts", a.middlewareAuth(getPostsForUserHandler.Handle))

	return mux
}
