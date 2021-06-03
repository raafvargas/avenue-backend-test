package storage

import (
	"context"
	"fmt"
	"io"

	cfg "github.com/avenuesec/file-upload-test/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	bucketName string
	client     *s3.Client
}

func NewS3Storage(cfg *cfg.Configuration) (*S3Storage, error) {
	endpointResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if cfg.S3Endpoint != "" {
			return aws.Endpoint{PartitionID: "aws", URL: cfg.S3Endpoint}, nil
		}

		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithEndpointResolver(endpointResolver),
	)

	if err != nil {
		return nil, err
	}

	storage := &S3Storage{
		bucketName: cfg.BucketName,
		client: s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.UsePathStyle = true
		}),
	}

	storage.ensureBucket(cfg.BucketName)

	return storage, nil
}

func (s *S3Storage) Copy(ctx context.Context, fileName string, reader io.Reader) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		Body:   reader,
	})

	return err
}

func (s *S3Storage) Read(ctx context.Context, fileName string, offset, limit int64) (io.Reader, error) {
	format := fmt.Sprintf("bytes=%d-%d", offset, offset+limit-1)

	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Range:  aws.String(format),
		Key:    aws.String(fileName),
		Bucket: aws.String(s.bucketName),
	})

	if err != nil {
		return nil, err
	}

	return output.Body, nil
}

func (s *S3Storage) ensureBucket(bucketName string) {
	s.client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
}
