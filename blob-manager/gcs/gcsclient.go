package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"log"
)

type GoogleStorageClient struct {
	client *storage.Client
	bucket string
}

func CreateGCSClient(ctx *context.Context, credentialFilePath, bucket string) *GoogleStorageClient {
	options := option.WithCredentialsFile(credentialFilePath)
	client, err := storage.NewClient(*ctx, options)
	if err != nil {
		log.Panicf("failed to create GCS Client , Reason: %s", err)
	}
	return &GoogleStorageClient{
		client: client,
		bucket: bucket,
	}
}
