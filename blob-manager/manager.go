package blobmanager

import "context"

type BlobManager struct {
	BlobStore BlobStore
}

func (b *BlobManager) Upload(ctx *context.Context, file File) error {
	return b.BlobStore.Upload(ctx, file)
}

func (b *BlobManager) Download(ctx *context.Context, fileName string) (File, error) {
	file, err := b.BlobStore.Download(ctx, fileName)
	return file, err
}

func (b *BlobManager) Delete(ctx *context.Context, fileName string) error {
	return b.BlobStore.Delete(ctx, fileName)
}
