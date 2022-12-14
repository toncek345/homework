package storage_test

import (
	"context"
	"io"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/spacelift-io/homework-object-storage/storage"
)

type minioMock struct {
	BucketExistsFn func(ctx context.Context, name string) (bool, error)
	MakeBucketFn   func(ctx context.Context, name string, opts minio.MakeBucketOptions) error
	GetObjectFn    func(ctx context.Context, bucket, id string, opts minio.GetObjectOptions) (*minio.Object, error)
	PutObjectFn    func(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
}

func (mm *minioMock) BucketExists(ctx context.Context, name string) (bool, error) {
	return mm.BucketExistsFn(ctx, name)
}

func (mm *minioMock) MakeBucket(ctx context.Context, name string, opts minio.MakeBucketOptions) error {
	return mm.MakeBucketFn(ctx, name, opts)
}

func (mm *minioMock) GetObject(ctx context.Context, bucket, id string, opts minio.GetObjectOptions) (*minio.Object, error) {
	return mm.GetObjectFn(ctx, bucket, id, opts)
}

func (mm *minioMock) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return mm.PutObjectFn(ctx, bucketName, objectName, reader, objectSize, opts)
}

func TestMinio(t *testing.T) {
	t.Run("initializes the bucket", func(t *testing.T) {
		storage.NewMinio(&minioMock{
			BucketExistsFn: func(ctx context.Context, name string) (bool, error) {
				if name != "testbucket" {
					t.Fatalf("wrong bucket name")
				}

				return false, nil
			},
			MakeBucketFn: func(ctx context.Context, name string, opts minio.MakeBucketOptions) error {
				if name != "testbucket" {
					t.Fatalf("wrong bucket name")
				}

				return nil
			},
		}, "testbucket")
	})

	// TODO: rest of the tests; get, put
}
