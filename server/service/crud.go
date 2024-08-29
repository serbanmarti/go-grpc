package service

import (
	"context"
	"fmt"
	"sync"

	"connectrpc.com/connect"
	"github.com/segmentio/ksuid"

	crudv1 "github.com/serbanmarti/go-grpc/proto_gen/crud/v1"
)

type CrudService struct {
	Mutex sync.RWMutex
	Data  map[string]string
}

func (s *CrudService) Create(ctx context.Context, req *connect.Request[crudv1.CreateRequest]) (*connect.Response[crudv1.CreateResponse], error) {
	// Create an ID for the new record
	id := ksuid.New().String()

	// Lock the mutex to ensure thread safety
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Store the record in memory
	s.Data[id] = req.Msg.Name

	return connect.NewResponse(&crudv1.CreateResponse{
		Id: id,
	}), nil
}

func (s *CrudService) Read(ctx context.Context, req *connect.Request[crudv1.ReadRequest]) (*connect.Response[crudv1.ReadResponse], error) {
	// Read lock the mutex to allow multiple readers
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	// Grab the record, if it exists
	name, ok := s.Data[req.Msg.Id]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("record not found"))
	}

	return connect.NewResponse(&crudv1.ReadResponse{
		Id:   req.Msg.Id,
		Name: name,
	}), nil
}

func (s *CrudService) Update(ctx context.Context, req *connect.Request[crudv1.UpdateRequest]) (*connect.Response[crudv1.UpdateResponse], error) {
	// Lock the mutex to ensure thread safety
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Check if the record exists
	_, ok := s.Data[req.Msg.Id]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("record not found"))
	}

	// Update the record
	s.Data[req.Msg.Id] = req.Msg.UpdatedName

	return connect.NewResponse(&crudv1.UpdateResponse{
		Id:   req.Msg.Id,
		Name: req.Msg.UpdatedName,
	}), nil
}

func (s *CrudService) Delete(ctx context.Context, req *connect.Request[crudv1.DeleteRequest]) (*connect.Response[crudv1.DeleteResponse], error) {
	// Lock the mutex to ensure thread safety
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Check if the record exists
	_, ok := s.Data[req.Msg.Id]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("record not found"))
	}

	// Delete the record
	delete(s.Data, req.Msg.Id)

	return connect.NewResponse(&crudv1.DeleteResponse{
		Id: req.Msg.Id,
	}), nil
}
