package word

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"wordganiser/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WordService interface {
	GetAllWords(context.Context) ([]models.Word, error)
	GetWordByName(name string, ctx context.Context) (*models.Word, error)
	GetWordById(name string, ctx context.Context) (*models.Word, error)
	AddWord(word *models.Word, ctx context.Context) (primitive.ObjectID, error)
	UpdateWord(id string, word *models.Word, ctx context.Context) (*models.Word, error)
}

type WordHandler struct {
	wordService WordService
}

func NewWordHandler(wordService WordService) *WordHandler {
	return &WordHandler{wordService: wordService}
}

func (wh *WordHandler) WordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		wh.getAllWordsHandler(w, r)
	case http.MethodPost:
		wh.postWordsHandler(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (wh *WordHandler) GetWordByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	if name == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	word, err := wh.wordService.GetWordByName(name, r.Context())

	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(word)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (wh *WordHandler) WordById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		wh.getWordById(id, w, r)
	case http.MethodPut:
		wh.updateWordById(id, w, r)
	case http.MethodPatch:
	case http.MethodDelete:
	default:
		w.Header().Set("Allow", "GET, PUT, PATCH, DELETE")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (wh *WordHandler) getWordById(id string, w http.ResponseWriter, r *http.Request) {
	word, err := wh.wordService.GetWordById(id, r.Context())

	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(word)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	w.Write(data)
}

func (wh *WordHandler) getAllWordsHandler(w http.ResponseWriter, r *http.Request) {
	words, err := wh.wordService.GetAllWords(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(words)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	w.Write(data)
}

func (wh *WordHandler) postWordsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var word models.Word
	err = json.Unmarshal(body, &word)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	id, err := wh.wordService.AddWord(&word, r.Context())

	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(id.String()))
}

func (wh *WordHandler) updateWordById(id string, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var word models.Word
	err = json.Unmarshal(body, &word)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	res, err := wh.wordService.UpdateWord(id, &word, r.Context())

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(res)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	w.Write(data)
}
