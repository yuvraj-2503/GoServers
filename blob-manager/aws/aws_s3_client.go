package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	client     *s3.S3
	config     *AWSConfig
}

func CreateS3Client(config *AWSConfig) *S3Client {
	sess := CreateSession(config)
	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)
	client := s3.New(sess)

	return &S3Client{
		uploader:   uploader,
		config:     config,
		downloader: downloader,
		client:     client,
	}
}

func CreateS3Session(sess *session.Session) *s3.S3 {
	s3Session := s3.New(sess)
	return s3Session
}

func CreateSession(awsConfig *AWSConfig) *session.Session {
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(awsConfig.Region),
			Credentials: credentials.NewStaticCredentials(
				awsConfig.AccessKeyID,
				awsConfig.AccessKeySecret,
				"",
			),
		},
	))
	return sess
}

type AWSConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Region          string
	BucketName      string
	UploadTimeout   int
	BaseURL         string
}
