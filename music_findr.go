package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"os"
	"time"
)

func apiIndex(w http.ResponseWriter, r *http.Request) {
	msg := struct {
		Status string `json:"status"`
	}{
		"OK",
	}

	render.JSON(w, r, msg)
}

func apiRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", apiIndex)

	return r
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	conn := fmt.Sprintf(":%s", port)

	router := chi.NewRouter()

	// Default middlewares
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.RequestID,
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.Recoverer,
		middleware.URLFormat,
		middleware.RedirectSlashes,
	)

	// Timeout
	router.Use(middleware.Timeout(60 * time.Second))

	// Mount API
	router.Mount("/api", apiRouter())

	log.Printf("Listening on: %s", conn)
	http.ListenAndServe(conn, router)
}
