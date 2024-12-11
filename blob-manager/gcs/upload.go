package gcs

import (
	blob_manager "blob-manager"
	common "blob-manager/common"
	"context"
	"io"
	"log"
)

func (g *GoogleStorageClient) Upload(ctx *context.Context, file blob_manager.File) error {
	defer close(file)
	wc := g.client.Bucket(g.bucket).Object(file.Name()).NewWriter(*ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return &common.UploadError{Message: err.Error()}
	}
	if err := wc.Close(); err != nil {
		return &common.UploadError{Message: err.Error()}
	}
	return nil
}

func close(file blob_manager.File) {
	err := file.Close()
	if err != nil {
		log.Fatalf("failed to close the file, reason: %s", err)
	}
}
