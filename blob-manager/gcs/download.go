package gcs

import (
	blob_manager "blob-manager"
	common "blob-manager/common"
	"cloud.google.com/go/storage"
	"context"
)

func (g *GoogleStorageClient) Download(ctx *context.Context, fileName string) (blob_manager.File, error) {
	rc, err := g.client.Bucket(g.bucket).Object(fileName).NewReader(*ctx)
	if err != nil {
		if storage.ErrObjectNotExist == err {
			return nil, &common.ObjectNotFound{ObjectId: fileName}
		}
		return nil, &common.DownloadError{
			Message: err.Error(),
		}
	}
	cloudStorageObject := &CloudStorageObject{
		FileName:      fileName,
		StorageReader: rc,
	}
	return cloudStorageObject, nil
}

type CloudStorageObject struct {
	FileName      string
	StorageReader *storage.Reader
}

func (o *CloudStorageObject) Type() string {
	return o.StorageReader.Attrs.ContentType
}

func (o *CloudStorageObject) Size() int64 {
	return o.StorageReader.Attrs.Size
}

func (o *CloudStorageObject) Read(p []byte) (n int, err error) {
	return o.StorageReader.Read(p)
}

func (o *CloudStorageObject) Name() string {
	return o.FileName
}

func (o *CloudStorageObject) Close() error {
	return o.StorageReader.Close()
}
