package main

import (
	"log"
	"net/http"

	"github.com/Aneeshie/loom/internal/agent"
	"github.com/Aneeshie/loom/internal/nlq"
	"github.com/Aneeshie/loom/internal/web/handlers"
	"github.com/Aneeshie/loom/internal/web/routes"
	"github.com/Aneeshie/loom/internal/web/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	client, err := agent.NewClient("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	parser := nlq.NewOllamaParser("llama3.2:3b")

	queryService := service.NewQueryService(parser, client)

	handler := handlers.NewHandler(queryService)

	routes.Register(r, handler)

	log.Fatal(http.ListenAndServe(":3000", r))
}
