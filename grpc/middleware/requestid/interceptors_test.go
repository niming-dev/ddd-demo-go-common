package requestid

import (
	"testing"

	grpc_testing "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

var (
	goodPing    = &pb_testproto.PingRequest{Value: "something", SleepTimeMs: 9999}
	reqId       = "3525ac55-cce8-4c4c-9520-6af78268050c"
	uuidBuilder = func() string {
		return reqId
	}
)

type ReqIdSuite struct {
	*grpc_testing.InterceptorTestSuite
}

func (s *ReqIdSuite) SetupTest() {
}

func (s *ReqIdSuite) TestPing_ServerGenRequestId() {
	opts := []Option{
		WithIdBuilder(func() string {
			return reqId + "_client"
		}),
	}
	s.Client = s.NewClient(grpc.WithChainUnaryInterceptor(UnaryClientInterceptor(opts...)))

	resp, err := s.Client.Ping(s.SimpleCtx(), goodPing)
	require.NoError(s.T(), err, "must not be an error on a successful call")

	tags := tagsFromJson(s.T(), resp.Value)
	require.Len(s.T(), tags, 1)
	assert.Equal(s.T(), tags["grpc.request.uuid"], reqId+"_client")
}

func (s *ReqIdSuite) TestPing_ServerGetClientReqId() {
	s.Client = s.NewClient()
	resp, err := s.Client.Ping(s.SimpleCtx(), goodPing)
	require.NoError(s.T(), err, "must not be an error on a successful call")

	tags := tagsFromJson(s.T(), resp.Value)
	require.Len(s.T(), tags, 1)
	assert.Equal(s.T(), tags["grpc.request.uuid"], reqId)
}

func TestReqIdSuiteSuite(t *testing.T) {
	opts := []Option{
		WithIdBuilder(uuidBuilder),
	}
	s := &ReqIdSuite{
		InterceptorTestSuite: &grpc_testing.InterceptorTestSuite{
			TestService: &tagPingBack{&grpc_testing.TestPingService{T: t}},
			ServerOpts: []grpc.ServerOption{
				grpc.StreamInterceptor(StreamServerInterceptor(opts...)),
				grpc.UnaryInterceptor(UnaryServerInterceptor(opts...)),
			},
		},
	}
	suite.Run(t, s)
}
