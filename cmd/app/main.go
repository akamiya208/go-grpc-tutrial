package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"

	"github.com/akamiya208/go-grpc-tutrial/internal/pkg/mysql"
	pb "github.com/akamiya208/go-grpc-tutrial/internal/pkg/proto"
	"github.com/akamiya208/go-grpc-tutrial/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mysqlClient, err := mysql.NewMySQLClient()
	if err != nil {
		slog.Error("failed to create mysql client")
		panic(err)
	}
	taskServer := server.NewTaskServer(mysqlClient)

	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("failed to listen: "+strconv.Itoa(port), slog.String("error", err.Error()))
		os.Exit(1)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterTaskServiceServer(s, taskServer)
	slog.Info("server listening at " + lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		slog.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
