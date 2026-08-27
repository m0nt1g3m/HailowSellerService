package s3_storage

import (
	"HailowSellerService/pkg/logging"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client   *s3.Client
	bucket   string
	region   string
	endpoint string
	logger   *logging.Logger
}

type Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	Endpoint        string
}

func NewS3Client(ctx context.Context, cfg Config, logger *logging.Logger) (*S3Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to load AWS SDK config: %w", err)
	}

	s3Opts := func(o *s3.Options) {
		if cfg.Endpoint != "" {
			endpoint := strings.TrimSuffix(cfg.Endpoint, "/")
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = false
		}
	}

	return &S3Client{
		client:   s3.NewFromConfig(awsCfg, s3Opts),
		bucket:   cfg.Bucket,
		region:   cfg.Region,
		endpoint: cfg.Endpoint,
		logger:   logger,
	}, nil
}

func (s *S3Client) Upload(ctx context.Context, folder string, filename string, file io.Reader, contentType string) (string, error) {
	objectKey := fmt.Sprintf("%s/%s", folder, filename)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("Failed to upload object to S3: %w", err)
	}

	var fileURL string
	if strings.Contains(s.endpoint, "supabase.co") {
		baseDomain := strings.TrimSuffix(strings.Split(s.endpoint, "/storage")[0], "/")
		fileURL = fmt.Sprintf("%s/storage/v1/object/public/%s/%s", baseDomain, s.bucket, objectKey)
	} else {
		fileURL = fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(s.endpoint, "/"), s.bucket, objectKey)
	}

	return fileURL, nil
}

func (s *S3Client) Delete(ctx context.Context, fileURL string) error {
	if fileURL == "" {
		return nil
	}

	objectKey := extractObjectKey(fileURL, s.bucket)
	if objectKey == "" {
		return fmt.Errorf("Unable to extract object key from URL: %s", fileURL)
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		s.logger.Warnf("Failed to delete object from S3. bucket=%s key=%s url=%s error=%v", s.bucket, objectKey, fileURL, err)
		return fmt.Errorf("Failed to delete object with key = %s from S3: %w", objectKey, err)
	}

	return nil
}

func extractObjectKey(fileURL string, bucket string) string {
	trimmed := strings.TrimSpace(fileURL)
	if trimmed == "" {
		return ""
	}

	if strings.Contains(trimmed, "/storage/v1/object/public/") {
		parts := strings.Split(trimmed, "/storage/v1/object/public/")
		if len(parts) < 2 {
			return ""
		}
		path := strings.TrimPrefix(parts[1], bucket+"/")
		return strings.TrimPrefix(path, "/")
	}

	for _, prefix := range []string{fmt.Sprintf("/%s/", bucket), fmt.Sprintf("%s.s3.", bucket)} {
		if strings.Contains(trimmed, prefix) {
			parts := strings.Split(trimmed, prefix)
			if len(parts) >= 2 {
				return strings.TrimPrefix(parts[1], "/")
			}
		}
	}

	return ""
}
