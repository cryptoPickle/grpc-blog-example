package blog

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type BlogItem struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string `bson:"author_id,omitempty"`
	Content string `bson:"content,omitempty"`
	Title string `bson:"title,omitempty"`
}