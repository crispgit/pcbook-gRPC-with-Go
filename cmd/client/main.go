package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/crispgit/pcbook/pb"
	"github.com/crispgit/pcbook/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	// connect to the server from client
	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't dial server", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	laptop := sample.NewLaptop()
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		statusCode, ok := status.FromError(err)
		if ok && statusCode.Code() == codes.AlreadyExists {
			log.Print("laptop already exists")
		} else {
			log.Fatal("can't create laptop: ", err)
		}
		return
	}

	log.Printf("create laptop with id: %s", res.Id)
}
