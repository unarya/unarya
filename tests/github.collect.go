package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/unarya/unarya/lib/proto/pb/collectorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("X:: Cannot connect to collector service: ", err)
	}
	defer conn.Close()
	client := collectorpb.NewCollectorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := collectorpb.GitRequest{
		Url:    "https://github.com/unarya/unarya",
		Branch: "master",
		Token:  "token",
	}

	res, err := client.CollectFromGit(ctx, &req)
	if err != nil {
		log.Fatal("X:: Cannot collect from git repository: ", err)
	}
	fmt.Println("Collected Git response:")
	fmt.Printf(" -> Path: %s", res.Path)
	fmt.Printf(" -> Message: %s", res.Message)
}
