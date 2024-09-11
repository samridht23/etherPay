package infra

import (
	"context"
	"fmt"
	"os"

	//"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type S3 struct {
	Client *s3.Client
}

func InitializeS3Client() (*S3, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_BUCKET_REGION")))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client: %v", err)
	}
	client := s3.NewFromConfig(cfg)
	zap.S().Info("S3 client initialized")

	s3Instance := &S3{
		Client: client,
	}
	//if err := s3Instance.PingBucket(os.Getenv("AWS_BUCKET_NAME")); err != nil {
	//	return nil, fmt.Errorf("failed to connect to S3 bucket: %v", err)
	//}
  return s3Instance,nil
}

// PingBucket checks if the connection to the specified S3 bucket is successful
//func (basic *S3) PingBucket(bucketName string) error {
//	_, err := basic.Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
//		Bucket: aws.String(bucketName),
//	})
//	if err != nil {
//		zap.S().Errorf("Failed to connect to S3 bucket %s: %v", bucketName, err)
//		return err
//	}
//	zap.S().Infof("Successfully connected to S3 bucket %s", bucketName)
//	return nil
//}

func (s3 *S3) NewClient() *s3.Client {
	return s3.Client
}
