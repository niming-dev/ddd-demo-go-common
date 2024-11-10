package nsq

import (
	"testing"

	"go.uber.org/zap"
)

func TestZapLogger_OutPut(t *testing.T) {
	zl, _ := zap.NewProduction()
	logger := ZapLogger{Logger: zl}
	_ = logger.Output(2, `ERR    3 [workflow.connector.poller_transfer/workflow.connector.poller_transfer] error querying nsqlookupd (http://nsqlookupd-http-4161-02.external-service.svc.cluster.local:4161/lookup?topic=workflow.connector.poller_transfer) - got response 404 Not Found "{\"message\":\"TOPIC_NOT_FOUND\"}"`)
}

func TestZapLogger_OutPut1(t *testing.T) {
	zl, _ := zap.NewProduction()
	logger := ZapLogger{Logger: zl}
	_ = logger.Output(2, `INF    2 [workflow.connector.request/workflow.connector] exiting lookupdLoop`)
}
