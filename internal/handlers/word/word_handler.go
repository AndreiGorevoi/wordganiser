package word

import (
	"context"
	"encoding/json"
	"net/http"
	"wordganiser/internal/models"
)

type WordService interface {
	GetAllWords(context.Context) ([]models.Word, error)
}

type WordHandler struct {
	wordService WordService
}

func NewWordHandler(wordService WordService) *WordHandler {
	return &WordHandler{wordService: wordService}
}

func (wh *WordHandler) AllWordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		words, err := wh.wordService.GetAllWords(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data, err := json.Marshal(words)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(data)
	default:
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
