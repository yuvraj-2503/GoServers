package blobmanager

import (
	"context"
	"io"
	"mime/multipart"
)

type File interface {
	io.Reader
	io.Closer
	Name() string
	Type() string
	Size() int64
}

type BlobStore interface {
	Upload(ctx *context.Context, file File) error
	Download(ctx *context.Context, fileName string) (File, error)
	Delete(ctx *context.Context, fileName string) error
}

type UploadableFile struct {
	fileId      string
	file        multipart.File
	fileSize    int64
	contentType string
}

func NewUploadableFile(
	fileId string,
	file multipart.File,
	size int64,
	contentType string) *UploadableFile {
	return &UploadableFile{
		fileId:      fileId,
		file:        file,
		fileSize:    size,
		contentType: contentType,
	}
}

func (u *UploadableFile) Read(p []byte) (n int, err error) {
	return u.file.Read(p)
}

func (u *UploadableFile) Name() string {
	return u.fileId
}

func (u *UploadableFile) Close() error {
	return u.file.Close()
}

func (u *UploadableFile) Size() int64 {
	return u.fileSize
}

func (u *UploadableFile) Type() string {
	return u.contentType
}
