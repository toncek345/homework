package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"hash/crc32"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioPool struct {
	instances []Storage
	crcTable  *crc32.Table
}

func NewMinioPool(dockerClient *client.Client, bucketName string) (Storage, error) {
	minioInstances := make([]Storage, 0, 10)

	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting container list: %w", err)
	}

	for _, container := range containers {
		if container.Image != "minio/minio" {
			continue
		}

		ip := ""
		for _, v := range container.NetworkSettings.Networks {
			ip = v.IPAddress
		}

		json, err := dockerClient.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			return nil, fmt.Errorf("inspecting container: %w", err)
		}

		secret, access := "", ""
		for _, v := range json.Config.Env {
			if strings.HasPrefix(v, "MINIO_ACCESS_KEY") {
				s := strings.Split(v, "=")
				access = s[1]
			}

			if strings.HasPrefix(v, "MINIO_SECRET_KEY") {
				s := strings.Split(v, "=")
				secret = s[1]
			}

			if secret != "" && access != "" {
				break
			}
		}

		mminio, err := minio.New(
			fmt.Sprintf("%s:9000", ip),
			&minio.Options{
				Secure: false,
				Creds:  credentials.NewStaticV4(access, secret, ""),
			})
		if err != nil {
			return nil, fmt.Errorf("connecting to minio container: %w", err)
		}

		m, err := NewMinio(mminio, bucketName)
		if err != nil {
			return nil, fmt.Errorf("creating minio: %w", err)
		}

		minioInstances = append(minioInstances, m)
	}

	return &MinioPool{
		instances: minioInstances,
		crcTable:  crc32.MakeTable(crc32.Castagnoli),
	}, nil
}

func NewTestMinioPool(s []Storage) MinioPool {
	return MinioPool{
		instances: s,
		crcTable:  crc32.MakeTable(crc32.Castagnoli),
	}
}

func (m *MinioPool) idHash(id string) int {
	return int(crc32.Checksum([]byte(id), m.crcTable)) % len(m.instances)
}

func (m *MinioPool) Get(ctx context.Context, id string) (io.ReadCloser, error) {
	return m.instances[m.idHash(id)].Get(ctx, id)
}

func (m *MinioPool) Put(ctx context.Context, id string, reader io.Reader, length int64) error {
	return m.instances[m.idHash(id)].Put(ctx, id, reader, length)
}
