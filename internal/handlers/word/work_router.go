package word

import (
	"net/http"
)

func RegisterWordRoutes(mux *http.ServeMux, wordHandler *WordHandler) {
	mux.HandleFunc("/words", wordHandler.AllWordsHandler)
}
