package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/twizar/teams/internal/adapters"
	"github.com/twizar/teams/internal/application/service"
	"github.com/twizar/teams/internal/ports"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	timeoutConnection                  = 10
	mongoConnURLEnvVar                 = "MONGO_CONN_URL"
	httpHeaderAccessControlAllowOrigin = "HTTP_HEADER_ACCESS_CONTROL_ALLOW_ORIGIN"
)

func main() {
	connURL, exists := os.LookupEnv(mongoConnURLEnvVar)
	if !exists {
		log.Panicf("required env `%s` is missing", mongoConnURLEnvVar)
	}

	accessControlAllowOrigin, exists := os.LookupEnv(httpHeaderAccessControlAllowOrigin)
	if !exists {
		log.Panicf("required env `%s` is missing", accessControlAllowOrigin)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutConnection*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURL))
	if err != nil {
		log.Panic("mongo connection error occurred")
	}

	repo := adapters.NewMongoDBTeamRepo(client)
	teams := service.NewTeams(repo)
	server := ports.NewHTTPServer(teams)
	r := ports.ConfigureRouter(server)
	adapter := gorillamux.New(r)
	handler := ports.NewLambdaHandler(adapter, accessControlAllowOrigin)
	lambda.Start(handler.Handle)
}
