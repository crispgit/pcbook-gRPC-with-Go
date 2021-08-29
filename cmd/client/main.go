package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/crispgit/pcbook/pb"
	"github.com/crispgit/pcbook/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createLaptop(laptopClient pb.LaptopServiceClient) {
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

func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
	log.Print("search filter: ", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("can't search laptop: ", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF{
			return
		}
		if err != nil {
			log.Fatal("can't receive response: ", err)
		}

		laptop := res.GetLaptop()
		log.Print("found: ", laptop.GetId())
	}
}

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

	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	searchLaptop(laptopClient, filter)
}
