package adapters

import (
	"context"
	"fmt"

	"github.com/twizar/common/pkg/dto"
	"github.com/twizar/teams/internal/domain/entity"
	"github.com/twizar/teams/internal/ports/converter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBName          = "twizar"
	TeamsCollection = "teams"
)

type MongoDBTeamRepo struct {
	client *mongo.Client
}

func NewMongoDBTeamRepo(client *mongo.Client) *MongoDBTeamRepo {
	return &MongoDBTeamRepo{client: client}
}

func (m MongoDBTeamRepo) All(ctx context.Context) ([]*entity.Team, error) {
	cursor, err := m.collection().Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("operation 'collection.Find' error: %w", err)
	}

	var result []dto.Team

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("operation 'cursor.All' error: %w", err)
	}

	return converter.DTOsToEntities(result), nil
}

func (m MongoDBTeamRepo) Filter(
	ctx context.Context,
	minRating float64,
	leagues []string,
	orderBy string,
	limit int,
) ([]*entity.Team, error) {
	ands := bson.D{
		{Key: "rating", Value: bson.D{{Key: "$gte", Value: minRating}}},
	}

	if len(leagues) > 0 {
		ands = append(ands, bson.E{Key: "league", Value: bson.D{{Key: "$in", Value: leagues}}})
	}

	filter := bson.D{
		{
			Key:   "$and",
			Value: bson.A{ands},
		},
	}

	cursor, err := m.collection().Find(
		ctx,
		filter,
		options.Find().SetSort(bson.D{{Key: orderBy, Value: -1}}).SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, fmt.Errorf("operation 'collection.Find' error: %w", err)
	}

	var result []dto.Team

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("operation 'cursor.All' error: %w", err)
	}

	return converter.DTOsToEntities(result), nil
}

func (m MongoDBTeamRepo) ByIDs(ctx context.Context, ids []string) ([]*entity.Team, error) {
	if ids == nil {
		ids = make([]string, 0)
	}

	cursor, err := m.collection().Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, fmt.Errorf("operation 'collection.Find' error: %w", err)
	}

	var result []dto.Team

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("operation 'cursor.All' error: %w", err)
	}

	return converter.DTOsToEntities(result), nil
}

func (m MongoDBTeamRepo) collection() *mongo.Collection {
	return m.client.Database(DBName).Collection(TeamsCollection)
}
