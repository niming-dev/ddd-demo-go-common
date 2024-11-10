package ufile

import (
	"context"
	"io"
	"time"

	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"

	"github.com/niming-dev/ddd-demo/go-common/log"
)

type Client struct {
	config ufsdk.Config
}

func NewClient(config ufsdk.Config) *Client {
	return &Client{
		config: config,
	}
}

func (a *Client) IOMutipartAsyncUpload(ctx context.Context, reader io.Reader, keyName string, mimeType string) error {
	req, err := ufsdk.NewFileRequest(&a.config, nil)
	defer hook(ctx, time.Now(), req)
	if err != nil {
		return err
	}
	return req.IOMutipartAsyncUpload(reader, keyName, mimeType)
}

func (a *Client) IOPut(ctx context.Context, reader io.Reader, keyName string, mimeType string) error {
	req, err := ufsdk.NewFileRequest(&a.config, nil)
	defer hook(ctx, time.Now(), req)
	if err != nil {
		return err
	}
	return req.IOPut(reader, keyName, mimeType)
}

func (a *Client) BucketName() string {
	return a.config.BucketName
}

func (a *Client) FileHost() string {
	return a.config.FileHost
}

func hook(ctx context.Context, start time.Time, req *ufsdk.UFileRequest) {
	if req == nil {
		return
	}
	log.Infof(ctx, "[%s] %+v", time.Since(start), string(req.DumpResponse(true)))
}
