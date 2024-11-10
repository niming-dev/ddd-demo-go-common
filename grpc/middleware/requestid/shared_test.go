package requestid

import (
	"context"
	"io"
	"testing"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	json "github.com/json-iterator/go"
)

func tagsToJson(value map[string]interface{}) string {
	str, _ := json.MarshalToString(value)
	return str
}

func tagsFromJson(t *testing.T, jstring string) map[string]interface{} {
	var msgMapTemplate interface{}
	err := json.UnmarshalFromString(jstring, &msgMapTemplate)
	if err != nil {
		t.Fatalf("failed unmarshaling tags from response %v", err)
	}
	return msgMapTemplate.(map[string]interface{})
}

type tagPingBack struct {
	pb_testproto.TestServiceServer
}

func (s *tagPingBack) Ping(ctx context.Context, ping *pb_testproto.PingRequest) (*pb_testproto.PingResponse, error) {
	return &pb_testproto.PingResponse{Value: tagsToJson(grpc_ctxtags.Extract(ctx).Values())}, nil
}

func (s *tagPingBack) PingError(ctx context.Context, ping *pb_testproto.PingRequest) (*pb_testproto.Empty, error) {
	return s.TestServiceServer.PingError(ctx, ping)
}

func (s *tagPingBack) PingList(ping *pb_testproto.PingRequest, stream pb_testproto.TestService_PingListServer) error {
	out := &pb_testproto.PingResponse{Value: tagsToJson(grpc_ctxtags.Extract(stream.Context()).Values())}
	return stream.Send(out)
}

func (s *tagPingBack) PingEmpty(ctx context.Context, empty *pb_testproto.Empty) (*pb_testproto.PingResponse, error) {
	return s.TestServiceServer.PingEmpty(ctx, empty)
}

func (s *tagPingBack) PingStream(stream pb_testproto.TestService_PingStreamServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		out := &pb_testproto.PingResponse{Value: tagsToJson(grpc_ctxtags.Extract(stream.Context()).Values())}
		err = stream.Send(out)
		if err != nil {
			return err
		}
	}
}
