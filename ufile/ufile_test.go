package ufile

import (
	"context"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	ufsdk "github.com/ufilesdk-dev/ufile-gosdk"

	mycontext "github.com/niming-dev/ddd-demo/go-common/context"
	"github.com/niming-dev/ddd-demo/go-common/log"
)

var (
	txt string = `
this is a test text.
`
	ufConfig ufsdk.Config = ufsdk.Config{
		PublicKey:       "TOKEN_fd9a8e69-72ac-4c27-b937-9ac1f207a4c5",
		PrivateKey:      "903b0b05-75d7-4d73-809a-e508601a0992",
		BucketHost:      "",
		BucketName:      "niming-test",
		FileHost:        "cn-bj.ufileos.com",
		VerifyUploadMD5: true,
	}
)

func TestIOPut(t *testing.T) {
	ctx, _ := mycontext.NewContext(context.Background(), log.NewFromLogrus(logrus.StandardLogger()))
	client := NewClient(ufConfig)
	reader := strings.NewReader(txt)
	keyName := "tmp/TestIOPut.txt"
	mimeType := "text/plain"
	err := client.IOPut(ctx, reader, keyName, mimeType)
	if err != nil {
		t.Error(err)
	}
}

func TestIOMutipartAsyncUpload(t *testing.T) {
	ctx, _ := mycontext.NewContext(context.Background(), log.NewFromLogrus(logrus.StandardLogger()))
	client := NewClient(ufConfig)
	reader := strings.NewReader(txt)
	keyName := "tmp/TestIOMutipartAsyncUpload.txt"
	mimeType := "text/plain"
	err := client.IOMutipartAsyncUpload(ctx, reader, keyName, mimeType)
	if err != nil {
		t.Error(err)
	}
}
