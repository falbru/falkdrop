package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Store struct {
	svc    *awsS3.PresignClient
	bucket string
}

func (store S3Store) NewUploadUrl(id string) (string, error) {
	req, err := store.svc.PresignPutObject(context.Background(), &awsS3.PutObjectInput{
		Bucket:      aws.String(store.bucket),
		Key:         aws.String(id),
		ContentType: aws.String("application/octet-stream"),
	}, func(opts *awsS3.PresignOptions) {
		opts.Expires = time.Minute * 15
	})

	return req.URL, err
}
