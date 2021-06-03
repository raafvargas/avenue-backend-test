package storage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/avenuesec/file-upload-test/pkg/config"
)

type GCSStorage struct {
	bucket *storage.BucketHandle
}

func NewGCSStorage(config *config.Configuration) (*GCSStorage, error) {
	client, err := storage.NewClient(context.Background())

	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(config.BucketName)

	storage := &GCSStorage{
		bucket: bucket,
	}

	return storage, nil
}

func (s *GCSStorage) Read(ctx context.Context, fileName string, offset, limit int64) (io.Reader, error) {
	object := s.bucket.Object(fileName)

	reader, err := object.NewRangeReader(ctx, offset, limit)

	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (s *GCSStorage) Copy(ctx context.Context, fileName string, reader io.Reader) error {
	object := s.bucket.Object(fileName)

	writer := object.NewWriter(ctx)

	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
