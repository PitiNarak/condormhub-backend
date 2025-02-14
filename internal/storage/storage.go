package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketType string

const (
	PublicBucket  BucketType = "public"
	PrivateBucket BucketType = "private"
)

type Config struct {
	BucketName        string `env:"BUCKET_NAME"`
	PrivateBucketName string `env:"PRIVATE_BUCKET_NAME"`
	AccountID         string `env:"ACCOUNT_ID"`
	AccessKeyID       string `env:"ACCESS_KEY_ID"`
	AccessKeySecret   string `env:"ACCESS_KEY_SECRET"`
	URL_PREFIX        string `env:"URL_PREFIX"`
}

type Storage struct {
	client        *s3.Client
	preSignClient *s3.PresignClient
	Config        Config
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
	client := NewClient(storageConfig)
	return &Storage{
		client:        client,
		Config:        storageConfig,
		preSignClient: s3.NewPresignClient(client),
	}
}

func (s *Storage) UploadFile(ctx context.Context, key string, contentType string, file io.Reader, bucket BucketType) error {
	var bucketName string
	if bucket == PublicBucket {
		bucketName = s.Config.BucketName
	} else {
		bucketName = s.Config.PrivateBucketName
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		ContentType: &contentType,
		Body:        file,
	})

	if err != nil {
		return err
	}

	return nil
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

func (s *Storage) CopyFile(ctx context.Context, sourceKey string, destKey string) error {
	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.Config.BucketName),
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.Config.BucketName, sourceKey)),
		Key:        aws.String(destKey),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) MoveFile(ctx context.Context, sourceKey string, destKey string) error {
	err := s.CopyFile(ctx, sourceKey, destKey)
	if err != nil {
		return err
	}

	err = s.DeleteFile(ctx, sourceKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetSignedUrl(ctx context.Context, key string, expires time.Duration) (string, error) {
	req, err := s.preSignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expires))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (s *Storage) GetPublicUrl(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.Config.URL_PREFIX, s.Config.BucketName, key)
}
