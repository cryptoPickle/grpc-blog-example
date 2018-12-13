package blog

import (
	"context"
	"fmt"
	"github.com/cryptoPickle/blog/Server/database"
	"github.com/cryptoPickle/blog/contract"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



type BlogService struct {
	Repository *database.MongoRepository
	Collection *mongo.Collection
}

func NewBlogService(r *database.MongoRepository, c *mongo.Collection) *BlogService {
	return &BlogService{
		Repository: r,
		Collection: c,
	}
}


func (b *BlogService) CreateBlog(c context.Context, req *contract.CreateBlogRequest) (*contract.CreateBlogResponse, error) {
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


func (b *BlogService) ReadBlog(c context.Context, req *contract.ReadBlogRequest) (*contract.ReadBlogResponse, error) {
	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot parse ID"))
	}

	data := &BlogItem{}
	filter := bson.M{"_id" : oid}

	err = b.Collection.FindOne(context.Background(), filter).Decode(data)

	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Cannot find blog with specified ID %v", err))
	}

	return &contract.ReadBlogResponse{
		Blog: &contract.Blog{
			Id: data.ID.Hex(),
			AuthorId: data.AuthorID,
			Contennt: data.Content,
			Title: data.Title,
		},
	}, nil

}

func (b *BlogService) UpdateBlog(ctx context.Context, req *contract.UpdateBlogRequest) (*contract.UpdateBlogResponse, error) {
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot parse ID"))
	}

	data := &BlogItem{
		Title: blog.Title,
		Content: blog.Contennt,
		AuthorID: blog.AuthorId,
	}
	resp := &BlogItem{}
	filter := bson.M{"_id" : oid}

	bd := bson.M{"$set" : data}
	var i options.ReturnDocument = 1

	err = b.Collection.FindOneAndUpdate(context.Background(), filter, bd, &options.FindOneAndUpdateOptions{
		ReturnDocument: &i,
	}).Decode(resp)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Not found %v", err))
	}

	return &contract.UpdateBlogResponse{
		Blog: &contract.Blog{
			Id: oid.Hex(),
			Title: resp.Title,
			Contennt: resp.Content,
			AuthorId: resp.AuthorID,
		},
	}, nil

}

func (b *BlogService) 	DeleteBlog(ctx context.Context, req *contract.DeleteBlogRequest) (*contract.DeleteBlogResponse, error) {
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Cannot parse ID"))
	}
	filter := bson.M{"_id" : oid}
	res, err := b.Collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot delete %v", err))
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find document %v", err))
	}

	return &contract.DeleteBlogResponse{
		Status: oid.Hex(),
	}, nil
}

func(b *BlogService) ListBlog(req *contract.ListBlogRequest, stream contract.BlogService_ListBlogServer) error {
	cursor, err := b.Collection.Find(context.Background(), nil)

	if err != nil {
		return  status.Error(codes.Internal, fmt.Sprintf("Unknown error %v", err))
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		data := &BlogItem{}
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Error While Decoding Data %v", err))
		}

		err = stream.Send(&contract.ListBlogResponse{Blog: &contract.Blog{
			Id: data.ID.Hex(),
			Title: data.Title,
			Contennt: data.Content,
			AuthorId: data.AuthorID,
		}})

		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Error while streaming %v", err))
		}
	}
	if err := cursor.Err(); err != nil {
		return  status.Error(codes.Internal, fmt.Sprintf("Unknown error %v", err))
	}
	return nil
}
