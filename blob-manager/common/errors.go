package gcs

import "fmt"

type ObjectNotFound struct {
	ObjectId string
}

func (e *ObjectNotFound) Error() string {
	return fmt.Sprintf("Object with id: %s doesn't exist", e.ObjectId)
}

type UploadError struct {
	Message string
}

type DownloadError struct {
	Message string
}

func (e *UploadError) Error() string {
	return "failed to upload file , reason:" + e.Message
}

func (e *DownloadError) Error() string {
	return "failed to download file, reason:" + e.Message
}
