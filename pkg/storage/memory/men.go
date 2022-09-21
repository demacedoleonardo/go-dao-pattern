package memory

import (
	"errors"
	"go-dao-pattern/pkg/context"
)

var (
	DataNotFoundErr = errors.New("memory data not found")
)

type StorageClient struct {
	m Memory
}

func (s *StorageClient) Get(ctx *context.Context, key string) (interface{}, error) {
	data, found := s.m[key]
	if !found {
		return nil, DataNotFoundErr
	}
	return data, nil
}

func (s *StorageClient) Save(ctx *context.Context, key string, value interface{}) error {
	s.m[key] = value
	return nil
}

func InitConnection() *StorageClient {
	client := new(StorageClient)
	client.m = make(Memory)
	return client
}
