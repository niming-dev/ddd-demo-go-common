package grpc

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	ufcontext "github.com/niming-dev/ddd-demo/go-common/context"
	"github.com/niming-dev/ddd-demo/go-common/log"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func setupContext() context.Context {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	return context.Background()
}

func TestDialContext(t *testing.T) {
	ctx := setupContext()

	dsn := os.Getenv("DSN")
	client, err := DialContext(ctx, dsn)

	t.Logf("client: %+v", client)
	t.Logf("err: %+v", err)
}

func TestClient(t *testing.T) {
	ctx := setupContext()
	ctx = ufcontext.WithLogFields(ctx, log.Fields{"a": "aaa"})

	dsn := os.Getenv("DSN")
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	client, err := DialContext(ctx, dsn)
	if err != nil {
		t.Error(err)
		return
	}
	c := pb.NewGreeterClient(client)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			name := "grpc_test"
			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
			if err != nil {
				t.Errorf("could not greet: %v", err)
			} else {
				t.Logf("Greeting: %s", r.GetMessage())
			}
		}
	}
}
