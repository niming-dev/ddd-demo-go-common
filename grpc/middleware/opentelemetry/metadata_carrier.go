package grpcopentelemetry

import (
	"google.golang.org/grpc/metadata"
)

type MetadataCarrier metadata.MD

// Get returns the value associated with the passed key.
func (mc MetadataCarrier) Get(key string) string {
	vs := metadata.MD(mc).Get(key)
	if len(vs) == 0 {
		return ""
	}

	return vs[0]
}

// Set stores the key-value pair.
func (mc MetadataCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (mc MetadataCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range mc {
		keys = append(keys, k)
	}
	return keys
}
