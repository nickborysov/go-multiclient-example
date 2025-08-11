package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpcHandler "github.com/nickborysov/go-multiclient-example/internal/grpc/handler"
	pb "github.com/nickborysov/go-multiclient-example/internal/grpc/proto"
	httpHandler "github.com/nickborysov/go-multiclient-example/internal/http/handler"
	mcpHandler "github.com/nickborysov/go-multiclient-example/internal/mcp/handler"
	"github.com/nickborysov/go-multiclient-example/internal/service"
	"google.golang.org/grpc"
)

func main() {
	srvc := service.New()

	httpRouter := httpHandler.NewRouter(srvc)
	grpcRouter := grpcHandler.NewRouter(srvc)
	mcpRouter := mcpHandler.NewRouter(srvc)

	httpRouter.Path("/mcp").
		Handler(mcpRouter.HTTPHandler())

	ctx, cancel := context.WithCancel(context.Background())

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh,
		syscall.SIGTERM, // terminate: stopped by `kill -9 PID`
		syscall.SIGINT,  // interrupt: stopped by Ctrl + C
	)

	go func(ctx context.Context) {
		port := ":8080"

		err := http.ListenAndServe(port, httpRouter)
		if err != nil {
			log.Fatal(err)
		}
	}(ctx)

	go func() {
		grpcServer := grpc.NewServer()
		pb.RegisterExampleServer(grpcServer, grpcRouter)
		ln, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal(err)
		}
		err = grpcServer.Serve(ln)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Server is running on port :8080 and gRPC on port :50051")
	<-exitCh
	cancel()
}
