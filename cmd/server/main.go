package main

import (
	"log"
	"net/http"
	"os"
	"wordganiser/configs"
	"wordganiser/internal/handlers/word"
	wordHandler "wordganiser/internal/handlers/word"
	wordRepository "wordganiser/internal/repositories/word"
	wordService "wordganiser/internal/services/word"
)

func main() {
	mux := http.NewServeMux()
	mongoDB, err := configs.ConnectMongoDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
		os.Exit(1)
	}
	wordRepository := wordRepository.NewMongoWordRepository(mongoDB)
	wordService := wordService.NewWordService(wordRepository)
	wordHandler := wordHandler.NewWordHandler(wordService)

	word.RegisterWordRoutes(mux, wordHandler)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	srv.ListenAndServe()
}
