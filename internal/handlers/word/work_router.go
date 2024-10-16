package word

import (
	"net/http"
)

func RegisterWordRoutes(mux *http.ServeMux, wordHandler *WordHandler) {
	mux.HandleFunc("/words", wordHandler.WordsHandler)
	mux.HandleFunc("/words/id/{id}", wordHandler.WordById)

	mux.HandleFunc("GET /words/name/{name}", wordHandler.GetWordByNameHandler)
}
