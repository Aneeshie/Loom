package main

import (
	"log"

	"github.com/Aneeshie/loom/internal/config"
	"github.com/Aneeshie/loom/internal/embeddings"
	"github.com/Aneeshie/loom/internal/grpc"
	"github.com/Aneeshie/loom/internal/storage"
)

func main() {

	cfg, err := config.LoadServerConfig("../../configs/server.yaml")
	if err != nil {
		log.Fatal(err)
	}
	store, err := storage.NewStore(cfg.Database.URL)
	if err != nil {
		log.Fatal(err)
	}

	defer store.CloseConnection()

	embedder := embeddings.NewOllamaEmbedder("nomic-embed-text")

	logService := grpc.NewLogService(store, embedder)

	server := grpc.NewServer(logService)
	server.Run()
}
