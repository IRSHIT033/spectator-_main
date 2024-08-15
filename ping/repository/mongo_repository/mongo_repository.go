package checkLogRepo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"spectator.main/domain"
	"spectator.main/internals/mongo"
)

type mongoRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	timeFormat     = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	collectionName = "checkLog"
)

func NewMongoRepository(DB mongo.Database) domain.CheckLogRepository {
	return &mongoRepository{DB, DB.Collection(collectionName)}
}

func (m *mongoRepository) InsertOne(ctx context.Context, data *domain.CheckLog) (*domain.CheckLog, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, data)
	if err != nil {
		return data, err
	}

	return data, nil
}
