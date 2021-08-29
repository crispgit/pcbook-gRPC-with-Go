package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/crispgit/pcbook/pb"
	"github.com/crispgit/pcbook/service"
	"google.golang.org/grpc"
)

func main() {
	// set port in command line
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("can't start the server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can't start the server: ", err)
	}
}
