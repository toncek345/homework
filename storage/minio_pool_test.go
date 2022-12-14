package storage

import (
	"context"
	"io"
	"testing"
)

func TestMinioPool(t *testing.T) {
	tests := []struct {
		name    string
		storage []Storage
		id      string
	}{
		{
			name: "write to first",
			id:   "a",
			storage: []Storage{
				&StorageMock{
					GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
						t := ctx.Value("test").(*testing.T)
						if id != "a" {
							t.Fatalf("received wrong id")
						}
						return nil, nil
					},
				},
				&StorageMock{
					GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
						t := ctx.Value("test").(*testing.T)
						t.Fatal("wrong container")
						return nil, nil
					},
				},
				&StorageMock{
					GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
						t := ctx.Value("test").(*testing.T)
						t.Fatal("wrong container")
						return nil, nil
					},
				},
			},
		},
		{
			name: "write to second",
			id:   "aa",
			storage: []Storage{
				&StorageMock{
					GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
						t := ctx.Value("test").(*testing.T)
						t.Fatal("wrong container")
						return nil, nil
					},
				},
				&StorageMock{
					GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
						t := ctx.Value("test").(*testing.T)
						if id != "aa" {
							t.Fatalf("received wrong id")
						}
						return nil, nil
					},
				},
				&StorageMock{
					GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
						t := ctx.Value("test").(*testing.T)
						t.Fatal("wrong container")
						return nil, nil
					},
				},
			},
		},
		{
			name: "write to third",
			id:   "aaaab",
			storage: []Storage{
				&StorageMock{
					GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
						t := ctx.Value("test").(*testing.T)
						t.Fatal("wrong container")
						return nil, nil
					},
				},
				&StorageMock{
					GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
						t := ctx.Value("test").(*testing.T)
						t.Fatal("wrong container")
						return nil, nil
					},
				},
				&StorageMock{
					GetFn: func(ctx context.Context, id string) (io.ReadCloser, error) {
						t := ctx.Value("test").(*testing.T)
						if id != "aaaab" {
							t.Fatalf("received wrong id")
						}
						return nil, nil
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool := NewTestMinioPool(test.storage)
			ctx := context.WithValue(context.Background(), "test", t)
			pool.Get(ctx, test.id)
		})
	}
}
