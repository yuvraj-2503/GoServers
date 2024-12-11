package gcs

import "context"

func (g *GoogleStorageClient) Delete(ctx *context.Context, fileName string) error {
	err := g.client.Bucket(g.bucket).Object(fileName).Delete(*ctx)
	return err
}
