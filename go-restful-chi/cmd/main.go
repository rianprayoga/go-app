package main

import (
	"example.com/chi-restful/pkg/datasource"
	"example.com/chi-restful/pkg/handlers/itemsHandler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	conn := datasource.NewPgClient("postgres://postgres:postgres@localhost:5432/go-app?sslmode=disable")

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// set allowed content-type in request
	r.Use(middleware.AllowContentType("application/json"))
	// set content-type in response
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Mount("/items", itemsHandler.New(conn))

	http.ListenAndServe(":8080", r)
}
