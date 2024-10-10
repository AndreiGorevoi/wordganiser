package word

import (
	"context"
	"time"
	"wordganiser/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoWordRepository struct {
	collection *mongo.Collection
}

func NewMongoWordRepository(db *mongo.Client) *MongoWordRepository {
	collection := db.Database("wordganiser").Collection("words")
	return &MongoWordRepository{collection: collection}
}

func (wr *MongoWordRepository) CreateWord(ctx context.Context, word *models.Word) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := wr.collection.InsertOne(ctx, word)

	if err != nil {
		return err
	}

	return nil
}

func (wr *MongoWordRepository) GetAllWords(ctx context.Context) ([]models.Word, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	crsr, err := wr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer crsr.Close(ctx)

	var words []models.Word

	for crsr.Next(ctx) {
		var word models.Word
		if err := crsr.Decode(&word); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	if err := crsr.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func (wr *MongoWordRepository) GetWordByName(ctx context.Context, name string) (*models.Word, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var word models.Word
	err := wr.collection.FindOne(ctx, bson.M{"word": name}).Decode(&word)

	if err != nil {
		return nil, err
	}

	return &word, nil
}
