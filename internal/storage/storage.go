package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	BucketName      string `env:"BUCKET_NAME"`
	AccountID       string `env:"ACCOUNT_ID"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	AccessKeySecret string `env:"ACCESS_KEY_SECRET"`
	URL_PREFIX      string `env:"URL_PREFIX"`
}

type Storage struct {
	client *s3.Client
	Config Config
}

func NewClient(storageConfig Config) *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(storageConfig.AccessKeyID, storageConfig.AccessKeySecret, "")),
		config.WithRegion("auto"),
	)

	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", storageConfig.AccountID))
	})

	return client
}

func NewStorage(storageConfig Config) *Storage {
	return &Storage{
		client: NewClient(storageConfig),
		Config: storageConfig,
	}
}

func (s *Storage) UploadFile(ctx context.Context, key string, contextType string, file io.Reader) (string, error) {
	output, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Config.BucketName),
		Key:         aws.String(key),
		ContentType: &contextType,
		Body:        file,
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.%s.r2.cf2.rackcdn.com/%s", s.Config.BucketName, s.Config.AccountID, url.PathEscape(key))

	fmt.Println(output)

	return url, nil
}

func (s *Storage) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(output.Body)

	return buf.Bytes(), nil
}

func (s *Storage) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return err
	}

	return nil
}
