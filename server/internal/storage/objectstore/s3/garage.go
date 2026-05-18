package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewGarageStore(bucket string) *S3Store {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return nil
	}

	svc := awsS3.NewFromConfig(cfg, func(o *awsS3.Options) {
		o.UsePathStyle = true
		o.Region = "garage"
	})

	return &S3Store{
		awsS3.NewPresignClient(svc),
		bucket,
	}
}
