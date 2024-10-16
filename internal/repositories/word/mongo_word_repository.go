package word

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wordganiser/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (wr *MongoWordRepository) GetWordByName(name string, ctx context.Context) (*models.Word, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var word models.Word
	err := wr.collection.FindOne(ctx, bson.M{"word": name}).Decode(&word)

	if err != nil {
		return nil, err
	}

	return &word, nil
}

func (wr *MongoWordRepository) GetWordById(id string, ctx context.Context) (*models.Word, error) {
	objid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var word models.Word
	err = wr.collection.FindOne(ctx, bson.M{"_id": objid}).Decode(&word)

	if err != nil {
		return nil, err
	}

	return &word, nil
}

func (wr *MongoWordRepository) AddWord(word *models.Word, ctx context.Context) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rs, err := wr.collection.InsertOne(ctx, word)

	if err != nil {
		return primitive.NilObjectID, err
	}

	insertedID, ok := rs.InsertedID.(primitive.ObjectID)

	if !ok {
		return primitive.NilObjectID, errors.New("failed to cast inserted ID to ObjectID")
	}

	return insertedID, nil
}

func (wr *MongoWordRepository) UpdateWord(id string, word *models.Word, ctx context.Context) (*models.Word, error) {
	exixting, err := wr.GetWordById(id, ctx)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	updateResult, err := wr.collection.UpdateByID(ctx, exixting.ID, bson.M{"$set": bson.M{"word": word.Word, "definition": word.Definition, "usage": word.Usage}})

	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, fmt.Errorf("no document found with ID: %s", id)
	}

	word.ID = exixting.ID
	return word, nil
}
