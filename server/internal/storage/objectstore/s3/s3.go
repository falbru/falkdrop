package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Store struct {
	client        *awsS3.Client
	presignClient *awsS3.PresignClient
	bucket        string
}

func NewS3Store(ctx context.Context, bucket string) (*S3Store, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return nil, err
	}

	svc := awsS3.NewFromConfig(cfg)

	return &S3Store{
		svc,
		awsS3.NewPresignClient(svc),
		bucket,
	}, nil
}

func (store S3Store) NewUploadUrl(ctx context.Context, id string) (string, error) {
	req, err := store.presignClient.PresignPutObject(ctx, &awsS3.PutObjectInput{
		Bucket:      aws.String(store.bucket),
		Key:         aws.String(id),
		ContentType: aws.String("application/octet-stream"),
	}, func(opts *awsS3.PresignOptions) {
		opts.Expires = time.Minute * 15
	})

	return req.URL, err
}

func (store S3Store) GetDownloadUrl(ctx context.Context, id string, filename string) (string, error) {
	req, err := store.presignClient.PresignGetObject(ctx, &awsS3.GetObjectInput{
		Bucket:                     aws.String(store.bucket),
		Key:                        aws.String(id),
		ResponseContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", filename)),
	}, func(opts *awsS3.PresignOptions) {
		opts.Expires = time.Minute * 15
	})

	return req.URL, err
}

func (store S3Store) CreateBucket(ctx context.Context) error {
	_, err := store.client.CreateBucket(ctx, &awsS3.CreateBucketInput{
		Bucket: aws.String(store.bucket),
		ACL:    types.BucketCannedACLPublicRead,
	})

	if err != nil {
		return err
	}

	return awsS3.NewBucketExistsWaiter(store.client).Wait(ctx, &awsS3.HeadBucketInput{Bucket: aws.String(store.bucket)}, time.Minute)
}
