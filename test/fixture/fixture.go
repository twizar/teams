package fixture

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/twizar/teams/internal/adapters"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoadFixtures(client *mongo.Client) error {
	collection := client.Database(adapters.DBName).Collection(adapters.TeamsCollection)

	byteValue, err := os.ReadFile("../../test/data/teams.json")
	if err != nil {
		return fmt.Errorf("couldn't read file: %w", err)
	}

	var teams []interface{}

	err = json.Unmarshal(byteValue, &teams)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal data: %w", err)
	}

	_, err = collection.InsertMany(context.Background(), teams)
	if err != nil {
		return fmt.Errorf("couldn't insert data: %w", err)
	}

	return nil
}

func ClearDB(client *mongo.Client) error {
	err := client.Database(adapters.DBName).Collection(adapters.TeamsCollection).Drop(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't drop collection: %w", err)
	}

	return nil
}
