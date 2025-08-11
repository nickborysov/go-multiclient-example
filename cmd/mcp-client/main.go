package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nickborysov/go-multiclient-example/internal/model"
)

const (
	uri = "http://localhost:8080/mcp"
)

func main() {
	ctx := context.Background()

	client := mcp.NewClient(&mcp.Implementation{Name: "example", Version: "v1.0.0"}, nil)
	cs, err := client.Connect(ctx, mcp.NewStreamableClientTransport(uri, nil))
	if err != nil {
		log.Fatal(err)
	}
	defer cs.Close()

	req := &mcp.CallToolParamsFor[any]{
		Meta: nil,
		Name: "sendInfo",
		Arguments: model.SendInfoRequest{
			FirstName: "go-client",
			LastName:  "",
			Email:     "",
			Phone:     "",
			Address:   "",
		},
	}

	for i := 0; i < 100; i++ {
		toolResult, err := cs.CallTool(ctx, req)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(toolResult)
		time.Sleep(time.Second)
	}

}
