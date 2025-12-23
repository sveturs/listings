package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
)

// Client handles MinIO/S3 operations for image storage
type Client struct {
	client        *minio.Client
	bucket        string
	publicBaseURL string
	logger        zerolog.Logger
}

// NewClient creates a new MinIO client
func NewClient(endpoint, accessKey, secretKey, bucket string, useSSL bool, logger zerolog.Logger) (*Client, error) {
	return NewClientWithPublicURL(endpoint, accessKey, secretKey, bucket, useSSL, "", logger)
}

// NewClientWithPublicURL creates a new MinIO client with a custom public base URL
func NewClientWithPublicURL(endpoint, accessKey, secretKey, bucket string, useSSL bool, publicBaseURL string, logger zerolog.Logger) (*Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Check if bucket exists, create if not
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		logger.Info().Str("bucket", bucket).Msg("created MinIO bucket")
	}

	// Construct public base URL if not provided
	if publicBaseURL == "" {
		protocol := "http"
		if useSSL {
			protocol = "https"
		}
		publicBaseURL = fmt.Sprintf("%s://%s/%s", protocol, endpoint, bucket)
	}

	logger.Info().
		Str("endpoint", endpoint).
		Str("bucket", bucket).
		Str("publicBaseURL", publicBaseURL).
		Bool("ssl", useSSL).
		Msg("MinIO client initialized")

	return &Client{
		client:        client,
		bucket:        bucket,
		publicBaseURL: publicBaseURL,
		logger:        logger.With().Str("component", "minio_client").Logger(),
	}, nil
}

// UploadImage uploads an image to MinIO
func (c *Client) UploadImage(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(
		ctx,
		c.bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		c.logger.Error().Err(err).Str("object", objectName).Msg("failed to upload image")
		return fmt.Errorf("failed to upload image: %w", err)
	}

	c.logger.Debug().Str("object", objectName).Int64("size", size).Msg("image uploaded successfully")
	return nil
}

// DownloadImage downloads an image from MinIO
func (c *Client) DownloadImage(ctx context.Context, objectName string) (io.ReadCloser, error) {
	object, err := c.client.GetObject(ctx, c.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		c.logger.Error().Err(err).Str("object", objectName).Msg("failed to get image")
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return object, nil
}

// DeleteImage deletes an image from MinIO
func (c *Client) DeleteImage(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(ctx, c.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		c.logger.Error().Err(err).Str("object", objectName).Msg("failed to delete image")
		return fmt.Errorf("failed to delete image: %w", err)
	}

	c.logger.Debug().Str("object", objectName).Msg("image deleted successfully")
	return nil
}

// GetPresignedURL generates a presigned URL for temporary access
func (c *Client) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, c.bucket, objectName, expiry, nil)
	if err != nil {
		c.logger.Error().Err(err).Str("object", objectName).Msg("failed to generate presigned URL")
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// GetPublicURL returns a permanent public URL for an object (no expiry).
// This should only be used when the bucket is configured with public read access.
func (c *Client) GetPublicURL(objectName string) string {
	return fmt.Sprintf("%s/%s", c.publicBaseURL, objectName)
}

// BucketExists checks if the bucket exists
func (c *Client) BucketExists(ctx context.Context) (bool, error) {
	exists, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return false, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	return exists, nil
}

// HealthCheck performs a health check on MinIO
func (c *Client) HealthCheck(ctx context.Context) error {
	_, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("MinIO health check failed: %w", err)
	}

	return nil
}
