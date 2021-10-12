package test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twizar/teams/internal/adapters"
	"github.com/twizar/teams/test/fixture"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TeamRepoHelper(t *testing.T) (repo *adapters.MongoDBTeamRepo, clean func()) {
	t.Helper()

	connURL, ok := os.LookupEnv("TEST_MONGO_CONN_URL")
	require.True(t, ok)

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	require.NoError(t, err)
	err = client.Ping(ctx, readpref.Primary())
	require.NoError(t, err)

	err = fixture.LoadFixtures(client)
	require.NoError(t, err)

	clean = func() {
		require.NoError(t, fixture.ClearDB(client))
	}

	return adapters.NewMongoDBTeamRepo(client), clean
}
