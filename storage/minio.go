package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	bucket string
	client minioClient
}

type minioClient interface {
	BucketExists(ctx context.Context, name string) (bool, error)
	MakeBucket(ctx context.Context, name string, opts minio.MakeBucketOptions) error
	GetObject(ctx context.Context, bucket, id string, opts minio.GetObjectOptions) (*minio.Object, error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
		opts minio.PutObjectOptions) (minio.UploadInfo, error)
}

func NewMinio(client minioClient, bucketName string) (Storage, error) {
	ok, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, fmt.Errorf("bucket check: %w", err)
	}

	if !ok {
		if err := client.MakeBucket(
			context.Background(),
			bucketName,
			minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("creating bucket: %w", err)
		}
	}

	return &MinioStorage{
		bucket: bucketName,
		client: client,
	}, nil
}

// Get returns the object from storage.
func (m *MinioStorage) Get(ctx context.Context, id string) (io.ReadCloser, error) {
	obj, err := m.client.GetObject(ctx, m.bucket, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("minio getting object: %w", err)
	}

	return obj, nil
}

// Put puts the object. If length of the object is unknown -1 will work but the Minio implementation
// will allocate a large amount of memory.
func (m *MinioStorage) Put(ctx context.Context, id string, reader io.Reader, length int64) error {
	_, err := m.client.PutObject(ctx, m.bucket, id, reader, length, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("minio putting the object: %w", err)
	}

	return nil
}
