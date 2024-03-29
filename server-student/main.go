package main

import (
	"log"
	"net"

	"github.com/tomassar/protobuffers-grpc-go/database"
	"github.com/tomassar/protobuffers-grpc-go/server"
	"github.com/tomassar/protobuffers-grpc-go/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", ":5060")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository("postgres://postgres:postgres@localhost:54321/student?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewStudentServer(repo)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on port 5060")
}
