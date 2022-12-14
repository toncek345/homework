package storage

import (
	"context"
	"io"
)

type StorageMock struct {
	GetFn func(ctx context.Context, id string) (io.ReadCloser, error)
	PutFn func(ctx context.Context, id string, reader io.Reader, length int64) error
}

func (m *StorageMock) Get(ctx context.Context, id string) (io.ReadCloser, error) {
	return m.GetFn(ctx, id)
}

func (m *StorageMock) Put(ctx context.Context, id string, reader io.Reader, length int64) error {
	return m.PutFn(ctx, id, reader, length)
}
