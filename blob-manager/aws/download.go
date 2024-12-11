package aws

import (
	blob_manager "blob-manager"
	common "blob-manager/common"
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
)

func (s *S3Client) Download(ctx *context.Context, fileName string) (blob_manager.File, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(fileName),
	}

	// Call the S3 GetObject API
	output, err := s.client.GetObject(input)
	if err != nil {
		log.Printf("failed to download file: %v", err)
		return nil, &common.DownloadError{
			Message: fmt.Sprintf("failed to download file with id: %s, reason: %s",
				fileName, err.Error())}
	}

	//defer output.Body.Close()

	// Read the content into a buffer
	buffer := new(bytes.Buffer)
	fileSize, err := io.Copy(buffer, output.Body)
	if err != nil {
		log.Printf("failed to read file content: %v", err)
		return nil, err
	}

	s3File := &S3File{
		content:  bytes.NewReader(buffer.Bytes()),
		fileName: fileName,
		fileType: output.ContentType,
		fileSize: fileSize,
	}

	return s3File, nil

}

type S3File struct {
	content     *bytes.Reader
	fileName    string
	fileType    *string
	fileSize    int64
	closeCalled bool
}

func (o *S3File) Type() string {
	return *o.fileType
}

func (o *S3File) Size() int64 {
	return o.fileSize
}

func (o *S3File) Read(p []byte) (n int, err error) {
	return o.content.Read(p)
}

func (o *S3File) Name() string {
	return o.fileName
}

func (o *S3File) Close() error {
	return nil
}
