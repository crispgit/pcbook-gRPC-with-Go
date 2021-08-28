package service

import (
	// "bytes"
	"context"
	"errors"
	// "io"
	"log"

	"github.com/google/uuid"
	"github.com/crispgit/pcbook/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server that provides laptop service
type LaptopServer struct {
	Store LaptopStore
}

// return a LaptopServer instance
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{store}
}

// CreateLaptop is a unary RPC to create a new laptop
// server *LaptopServer is the receiver
func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		// check if it is a valid uuid
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid uuid: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	// save the laptop to store
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}

		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("saved laptop with id: %s", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}
