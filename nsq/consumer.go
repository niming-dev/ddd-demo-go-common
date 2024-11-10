package nsq

import (
	"github.com/nsqio/go-nsq"
)

func NewConsumer(topic string, channel string, config *nsq.Config) (*nsq.Consumer, error) {
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}
	consumer.SetLogger(&cunsumerLogger{}, nsq.LogLevelDebug)

	return consumer, nil
}
