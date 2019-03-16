package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
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

	log.Printf("Listening on port %d ...", 4000)
	http.ListenAndServe(":4000", router)
}
