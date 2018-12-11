package database

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)



type MongoRepository struct {
	client *mongo.Client
	collection *mongo.Collection
}

func NewMongo (c *mongo.Client) (*MongoRepository, *mongo.Collection) {
	mc := &MongoRepository{client: c}
	collection := mc.CreateCollection("gotest", "blog")
	mc.collection = collection
	return mc, collection
}

func (m *MongoRepository) CreateCollection(dbname , cname string) *mongo.Collection {
	collection := m.client.Database(dbname).Collection(cname)
	return collection
}