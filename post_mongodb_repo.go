package main

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type PostMongoRepo struct {
	client *mongo.Client
}
