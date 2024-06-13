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

	mux.Handle("POST /users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := &PostUserHandler{
			DB: a.Cfg.DB,
		}
		handler.Handle(w, r)
	}))

    mux.Handle("GET /user", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        handler := &GetUserHandler{
            DB: a.Cfg.DB,
        }
        handler.Handle(w, r)
    }))

	return mux
}
