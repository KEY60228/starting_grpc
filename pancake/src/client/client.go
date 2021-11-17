package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"pancake/client/gen/api"
)

var client = getClient()
var md = metadata.New(map[string]string{"authorization": "bearer hi/mi/tsu"})

func getClient() api.PancakeBakerServiceClient {
	address := "api:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can not connect: %v", err)
	}

	return api.NewPancakeBakerServiceClient(conn)
}

func bakePancake(menu api.Pancake_Menu) (*api.BakeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx = metadata.NewOutgoingContext(ctx, md)
	defer cancel()

	req := &api.BakeRequest{
		Menu: menu,
	}
	r, err := client.Bake(ctx, req)

	return r, err
}

func getReport() (*api.ReportResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx = metadata.NewOutgoingContext(ctx, md)
	defer cancel()

	req := &api.ReportRequest{}
	r, err := client.Report(ctx, req)

	return r, err
}

func main() {
	res, _ := bakePancake(api.Pancake_CLASSIC)
	fmt.Printf("Chef %s served %s pancake.\n", res.Pancake.ChefName, res.Pancake.Menu)

	res, _ = bakePancake(api.Pancake_MIX_BERRY)
	fmt.Printf("Chef %s served %s pancake.\n", res.Pancake.ChefName, res.Pancake.Menu)

	res, _ = bakePancake(api.Pancake_CLASSIC)
	fmt.Printf("Chef %s served %s pancake.\n", res.Pancake.ChefName, res.Pancake.Menu)

	reportRes, _ := getReport()
	for _, bakeCount := range reportRes.Report.BakeCounts {
		fmt.Println(bakeCount.Menu, bakeCount.Count)
	}
}
