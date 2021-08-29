package service

import (
	// "context"
	"errors"
	"fmt"
	// "log"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/crispgit/pcbook/pb"
)

// ErrAlreadyExists is returned when a record with the same ID already exists in the store
var ErrAlreadyExists = errors.New("record already exists")

// Laptop store is the interface to store laptop
type LaptopStore interface {
	// Save the laptop to the store
	Save (laptop *pb.Laptop) error
	// Find a laptop by id
	Find (id string) (*pb.Laptop, error)
}

// Store laptops in memory
type InMemoryLaptopStore struct {
	// read write mutex to handle concurrency
	mutex sync.RWMutex
	// laptop.id -> laptop object
	data map[string]*pb.Laptop
}

// returns a new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// Saves the laptop to the store
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	// require a write lock before adding new object
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[other.Id] = other
	return nil
}

// Find a laptop by id
func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	return deepCopy(laptop)
}


func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}


/* Implement the db store in the future
type DBLaptopStore struct {

}
*/
