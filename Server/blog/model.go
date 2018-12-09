package blog

import "gopkg.in/mgo.v2/bson"

type BlogItem struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	AuthorID string `bson:"author_id,omitempty"`
	Content string `bson:"content,omitempty"`
	Title string `bson:"title,omitempty"`
}