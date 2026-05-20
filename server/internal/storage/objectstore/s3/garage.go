package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewGarageStore(ctx context.Context, bucket string) (*S3Store, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return nil, err
	}

	svc := awsS3.NewFromConfig(cfg, func(o *awsS3.Options) {
		o.UsePathStyle = true
		o.Region = "garage"
	})

	return &S3Store{
		svc,
		awsS3.NewPresignClient(svc),
		bucket,
	}, nil
}
