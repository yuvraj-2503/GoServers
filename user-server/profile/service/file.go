package service

import "io"

type UploadableFile struct {
	fileId      string
	file        io.ReadCloser
	size        int64
	contentType string
}

func NewUploadableFile(
	fileId string,
	file io.ReadCloser,
	fileSize int64,
	contentType string) *UploadableFile {
	return &UploadableFile{
		fileId:      fileId,
		file:        file,
		size:        fileSize,
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
	return u.size
}

func (u *UploadableFile) Type() string {
	return u.contentType
}
