package storage

import (
	"context"
	"io"
)

var _ Storage = (*MinioStorage)(nil)
var _ Storage = (*StorageMock)(nil)
var _ Storage = (*StoragePool)(nil)

type Storage interface {
	Get(ctx context.Context, id string) (io.ReadCloser, error)
	Put(ctx context.Context, id string, reader io.Reader, length int64) error
}
