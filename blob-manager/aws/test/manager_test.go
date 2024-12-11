package test

import (
	"blob-manager"
	"blob-manager/aws"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"
	"testing"
)

func TestBlobManager_Delete(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		BlobStore blobmanager.BlobStore
	}
	type args struct {
		ctx      *context.Context
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Delete",
			fields: fields{
				BlobStore: aws.CreateS3Client(&aws.AWSConfig{
					AccessKeyID:     "AKIAXSWKI2KJD4K6QQDV",
					AccessKeySecret: "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",
					Region:          "ap-south-1",
					BucketName:      "yuvrajusersprofile",
					UploadTimeout:   0,
					BaseURL:         "https://yuvrajusersprofile.s3.ap-south-1.amazonaws.com",
				}),
			},
			args: args{
				ctx:      &ctx,
				fileName: "674c9d65361369da5c4710d6",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &blobmanager.BlobManager{
				BlobStore: tt.fields.BlobStore,
			}
			if err := b.Delete(tt.args.ctx, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlobManager_Download(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		BlobStore blobmanager.BlobStore
	}
	type args struct {
		ctx      *context.Context
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    blobmanager.File
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Download",
			fields: fields{
				BlobStore: aws.CreateS3Client(&aws.AWSConfig{
					AccessKeyID:     "AKIAXSWKI2KJD4K6QQDV",
					AccessKeySecret: "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",
					Region:          "ap-south-1",
					BucketName:      "yuvrajusersprofile",
					UploadTimeout:   0,
					BaseURL:         "https://yuvrajusersprofile.s3.ap-south-1.amazonaws.com",
				}),
			},
			args: args{
				ctx:      &ctx,
				fileName: "674c9d65361369da5c4710d64",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &blobmanager.BlobManager{
				BlobStore: tt.fields.BlobStore,
			}
			got, err := b.Download(tt.args.ctx, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fileContent, err := io.ReadAll(got)
			if err != nil {
				t.Errorf("ReadAll() error = %v", err)
			}

			base64 := base64.StdEncoding.EncodeToString(fileContent)
			fmt.Println("Base64: " + base64)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Download() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlobManager_Upload(t *testing.T) {
	var ctx = context.Background()
	type fields struct {
		BlobStore blobmanager.BlobStore
	}
	type args struct {
		ctx  *context.Context
		file blobmanager.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Upload",
			fields: fields{
				BlobStore: aws.CreateS3Client(&aws.AWSConfig{
					AccessKeyID:     "AKIAXSWKI2KJD4K6QQDV",
					AccessKeySecret: "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",
					Region:          "ap-south-1",
					BucketName:      "yuvrajusersprofile",
					UploadTimeout:   0,
					BaseURL:         "https://yuvrajusersprofile.s3.ap-south-1.amazonaws.com",
				}),
			},
			args: args{
				ctx:  &ctx,
				file: getFile("674c9d65361369da5c4710d6", base64Image),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &blobmanager.BlobManager{
				BlobStore: tt.fields.BlobStore,
			}
			if err := b.Upload(tt.args.ctx, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getFile(fileId, fileContent string) blobmanager.File {
	decodedBytes, _ := base64.StdEncoding.DecodeString(fileContent)
	return blobmanager.NewUploadableFile(fileId, getReadCloserFromByteArray(decodedBytes),
		int64(len(decodedBytes)), "image/*")
}

// A struct to wrap bytes.Reader and implement io.Closer
type reader struct {
	*bytes.Reader
}

// Implement the Close method to satisfy the io.Closer interface
func (rc *reader) Close() error {
	return nil // No resources to release, so return nil
}

// getReadCloserFromByteArray converts a byte array to multipart.File
func getReadCloserFromByteArray(data []byte) multipart.File {
	return &reader{Reader: bytes.NewReader(data)}
}
