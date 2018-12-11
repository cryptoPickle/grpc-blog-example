package blog

import (
	"context"
	"fmt"
	"github.com/cryptoPickle/blog/Server/database"
	"github.com/cryptoPickle/blog/contract"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



type BlogServiceServer struct {
	Repository *database.MongoRepository
	Collection *mongo.Collection
}

func NewBlogService(r *database.MongoRepository, c *mongo.Collection) *BlogServiceServer {
	return &BlogServiceServer{
		Repository: r,
		Collection: c,
	}
}


func (b *BlogServiceServer) CreateBlog(c context.Context, req *contract.CreateBlogRequest) (*contract.CreateBlogResponse, error) {
	blog := req.GetBlog()
	data := BlogItem{
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContennt(),
	}
	res, err := b.Collection.InsertOne(context.Background(), data)

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error %v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to OID %v", oid))
	}


	return &contract.CreateBlogResponse{
		Blog: &contract.Blog{
			Id: oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Contennt: blog.GetContennt(),
		},
	}, nil
}

