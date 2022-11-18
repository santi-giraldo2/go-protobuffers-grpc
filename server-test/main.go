package main

import (
	"log"
	"net"

	"enerbit.com/go/grpc/database"
	"enerbit.com/go/grpc/server"
	"enerbit.com/go/grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", ":5070")

	if err != nil {
		log.Fatal(err)
	}

	defer list.Close()
	repo, err := database.NewPostgresRepository("postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	server := server.NewTestServer(repo)
	s := grpc.NewServer()

	testpb.RegisterTestServiceServer(s, server)
	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
