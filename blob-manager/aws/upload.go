package aws

import (
	blobmanager "blob-manager"
	common "blob-manager/common"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (s *S3Client) Upload(ctx *context.Context, file blobmanager.File) error {
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(file.Name()),
		Body:   file,
	})

	if err != nil {
		fmt.Printf("failed to upload object, %v\n", err)
		return &common.UploadError{
			Message: fmt.Sprintf("failed to upload object: %s,reason: %v\n",
				file.Name(), err.Error())}
	}

	fmt.Printf("Successfully uploaded %q to %q\n", file.Name(), s.config.BucketName)
	return nil
}
