package main

import (
	"context"
	"log"
	"time"

	"github.com/unarya/unarya/lib/proto/pb/collectorpb"
	"github.com/unarya/unarya/tests/collector/client"
)

func main() {
	c, err := client.NewCollectorClient("localhost:50051", "", 5)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	res, err := c.Client.Ready(ctx, &collectorpb.Empty{})
	cancel()
	log.Println("[RESPONSE]::" + res.Status)
}
