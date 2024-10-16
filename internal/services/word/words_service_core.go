package word

import (
	"context"
	"wordganiser/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WordRepository interface {
	GetAllWords(context.Context) ([]models.Word, error)
	GetWordByName(name string, ctx context.Context) (*models.Word, error)
	GetWordById(id string, ctx context.Context) (*models.Word, error)
	AddWord(word *models.Word, ctx context.Context) (primitive.ObjectID, error)
	UpdateWord(id string, word *models.Word, ctx context.Context) (*models.Word, error)
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

func (ws *WordServiceCore) GetWordByName(name string, ctx context.Context) (*models.Word, error) {
	return ws.repository.GetWordByName(name, ctx)
}

func (ws *WordServiceCore) GetWordById(name string, ctx context.Context) (*models.Word, error) {
	return ws.repository.GetWordById(name, ctx)
}

func (ws *WordServiceCore) AddWord(word *models.Word, ctx context.Context) (primitive.ObjectID, error) {
	return ws.repository.AddWord(word, ctx)
}

func (ws *WordServiceCore) UpdateWord(id string, word *models.Word, ctx context.Context) (*models.Word, error) {
	return ws.repository.UpdateWord(id, word, ctx)
}
