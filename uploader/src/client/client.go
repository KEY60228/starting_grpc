package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	"image/client/gen/api"
)

var client = getClient()

func getClient() api.ImageUploadServiceClient {
	address := "api:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return api.NewImageUploadServiceClient(conn)
}

func uploadImage(filepath string) (*api.ImageUploadResponse, error) {
	fp, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, _ := client.Upload(ctx)

	filemeta := &api.ImageUploadRequest_FileMeta_{
		FileMeta: &api.ImageUploadRequest_FileMeta{
			Filename: filepath,
		},
	}
	req := &api.ImageUploadRequest{File: filemeta}
	stream.Send(req)
	fmt.Println("sent", filepath)

	buf := make([]byte, 100*1024)
	for {
		n, _ := fp.Read(buf)
		if n == 0 {
			break
		}
		data := &api.ImageUploadRequest_Data{Data: buf[:n]}
		req := &api.ImageUploadRequest{File: data}
		stream.Send(req)
		fmt.Println("sent", n)
	}

	r, err := stream.CloseAndRecv()
	return r, err
}

func main() {
	filepath := flag.String("f", "default.jpeg", "Image path")
	flag.Parse()
	r, _ := uploadImage(*filepath)
	fmt.Println(r)
}
