package word

import (
	"context"
	"wordganiser/internal/models"
)

type WordRepository interface {
	GetAllWords(context.Context) ([]models.Word, error)
}

type WordServiceCore struct {
	repository WordRepository
}

func NewWordService(repository WordRepository) *WordServiceCore {
	return &WordServiceCore{repository: repository}
}

func (ws *WordServiceCore) GetAllWords(ctx context.Context) ([]models.Word, error) {
	return ws.repository.GetAllWords(ctx)
}
