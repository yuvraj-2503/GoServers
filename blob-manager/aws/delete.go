package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

func (s *S3Client) Delete(ctx *context.Context, fileName string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		log.Printf("Error deleting file %s: %s", fileName, err)
		return err
	}

	log.Printf("Object %s successfully deleted from bucket %s", fileName, s.config.BucketName)
	return nil
}
