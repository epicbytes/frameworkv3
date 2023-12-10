package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CounterModel struct {
	collection *mongo.Collection
}

func (s *CounterModel) NextID(ctx context.Context, collection string) (uint32, error) {
	res := s.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": collection},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true),
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	if err := res.Err(); err != nil {
		return 0, errors.Errorf("failed to find one and update: %w", err)
	}

	var counter struct {
		ID    string `bson:"_id"`
		Value uint32 `bson:"value"`
	}
	if err := res.Decode(&counter); err != nil {
		return 0, errors.Errorf("failed to decode counter: %w", err)
	}

	return counter.Value, nil
}
