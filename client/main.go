package main

import (
	"context"
	"io"
	"log"
	"sync"
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
	DoBidirectionalStreaming(c)
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

func DoServerStreaming(c testpb.TestServiceClient) {
	req := &testpb.GetStudentsPerTestRequest{
		TestId: "t3",
	}

	stream, err := c.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetStudentsPerTest RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error while receiving response: %v", err)
			break
		}

		log.Printf("Response from GetStudentsPerTest: %v", res)
	}
}

func DoBidirectionalStreaming(c testpb.TestServiceClient) {
	answer := testpb.TakeTestRequest{
		Answer: "2",
	}

	numberOfQuestions := 4

	wg := sync.WaitGroup{}
	wg.Add(1)

	stream, err := c.TakeTest(context.Background())
	if err != nil {
		log.Fatalf("error while opening stream: %v", err)
	}

	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			log.Printf("Sending answer: %v", &answer)
			stream.Send(&answer)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error while receiving response: %v", err)
				break
			}

			log.Printf("Response from TakeTest: %v", res)
		}

		wg.Done()
	}()

	wg.Wait()
}
