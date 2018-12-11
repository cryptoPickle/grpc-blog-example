package main

import (
	"context"
	"github.com/cryptoPickle/blog/contract"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	c := contract.NewBlogServiceClient(cc)
	CreateBlog(c)
}

func CreateBlog(b contract.BlogServiceClient) {
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
}