package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

func getConn() *sql.DB {
	dbname := os.Getenv("PG_DATABASE")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s port=%s host=%s sslmode=disable",
		user,
		dbname,
		password,
		port,
		host,
	)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func apiIndex(w http.ResponseWriter, r *http.Request) {
	var version string

	db := getConn()

	err := db.QueryRow("SHOW server_version").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	msg := struct {
		Status    string `json:"status"`
		DBVersion string `json:"dbversion"`
	}{
		"OK",
		version,
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
