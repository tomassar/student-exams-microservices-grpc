package main

import (
	"context"
	"log"
	"time"

	"github.com/tomassar/protobuffers-grpc-go/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := testpb.NewTestServiceClient(cc)
	DoClientStreaming(c)
}

func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "t3",
	}

	res, err := c.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetTest RPC: %v", err)
	}

	log.Printf("Response from GetTest: %v", res)
}

func DoClientStreaming(c testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "q4t3",
			Question: "What is your name?",
			Answer:   "My name is Tomassar",
			TestId:   "t3",
		},
		{
			Id:       "q5t3",
			Question: "What is your favourite programming language?",
			Answer:   "My favourite programming language is Go",
			TestId:   "t3",
		},
		{
			Id:       "q6t3",
			Question: "What is your favourite IDE?",
			Answer:   "My favourite IDE is Goland",
			TestId:   "t3",
		},
	}

	stream, err := c.SetQuestions(context.Background())
	if err != nil {
		log.Fatalf("error while opening stream: %v", err)
	}

	ticker := time.NewTicker(2 * time.Second)
	for _, question := range questions {
		log.Printf("Sending question: %v", question)
		stream.Send(question)
		<-ticker.C
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
	}

	log.Printf("SetQuestions Response: %v", res)
}
