package main

import (
	"context"
	"github.com/cryptoPickle/blog/contract"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"io"
)

func main() {
	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	c := contract.NewBlogServiceClient(cc)
	id := CreateBlog(c)
	ReadBlog(c, id)
	UpdateBlog(c, id)
	//DeleteBlog(c, id)
	ListBlog(c)
}

func ListBlog(b contract.BlogServiceClient) {
	stream, err := b.ListBlog(context.Background(), &contract.ListBlogRequest{})

	if err != nil {
		log.Fatalf("Error on request %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error on recieving %v", err)
		}

		log.Infof("Streamed ===============> %v", res)
	}

}

func DeleteBlog(b contract.BlogServiceClient, id string) {
	resp, err := b.DeleteBlog(context.Background(), &contract.DeleteBlogRequest{Id: id})

	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Blog deleted %v", resp)
}
func UpdateBlog(b contract.BlogServiceClient, id string) {
	update :=  &contract.UpdateBlogRequest{
		Blog: &contract.Blog{
			Id: id,
			Title: "Updated title",
			Contennt: "Updated content",
			AuthorId: "1231",
		},
	}
	resp, err := b.UpdateBlog(context.Background(),update)

	if err != nil {
		log.Fatal(err)
	}

	log.Info("============================================")
	log.Warn(resp)
}

func ReadBlog(b contract.BlogServiceClient, Id string) {
	resp, err := b.ReadBlog(context.Background(), &contract.ReadBlogRequest{BlogId: Id})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(resp)
}

func CreateBlog(b contract.BlogServiceClient) string {
	log.Info("Creating Blog....")
	blog :=  &contract.Blog{
		AuthorId: "Test",
		Title: "My first blog",
		Contennt: "first blog",
	}
	resp , err := b.CreateBlog(context.Background(), &contract.CreateBlogRequest{Blog:blog})

	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Blog has been created %v", resp)
	return resp.Blog.Id
}