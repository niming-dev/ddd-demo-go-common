package context

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/niming-dev/ddd-demo/go-common/log"
	"github.com/sirupsen/logrus"
)

func setupTest() context.Context {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	logger.SetFormatter(new(logrus.JSONFormatter))
	logger.SetLevel(logrus.DebugLevel)

	ctx, err := NewContext(context.Background(), log.NewFromLogrus(logger))
	if err != nil {
		panic(err)
	}
	return ctx
}

func TestDebug(t *testing.T) {
	ctx := setupTest()
	t.Logf("debug %+v", ctx)
}
