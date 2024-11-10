package ufgateway

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestExtractGatewayData(t *testing.T) {
	tables := []struct {
		name   string
		in     map[string][]string
		expect GatewayData
	}{
		{
			name: "full",
			in: map[string][]string{
				GWProjectIdKey:       {"prj_id_test1"},
				GWProjectNameKey:     {"prj_name_test1"},
				GWAccountIdentityKey: {"account_id_test1"},
				GWAccountNameKey:     {"account_name_test1"},
			},
			expect: GatewayData{
				ProjectId:       "prj_id_test1",
				ProjectName:     "prj_name_test1",
				AccountIdentity: "account_id_test1",
				AccountName:     "account_name_test1",
			},
		},
		{
			name: "omit_prj",
			in: map[string][]string{
				GWProjectIdKey:       {},
				GWProjectNameKey:     {},
				GWAccountIdentityKey: {"account_id_test1"},
				GWAccountNameKey:     {"account_name_test1"},
			},
			expect: GatewayData{
				ProjectId:       "",
				ProjectName:     "",
				AccountIdentity: "account_id_test1",
				AccountName:     "account_name_test1",
			},
		},
	}

	for _, v := range tables {
		t.Run(v.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), v.in)
			gwData, ok := ExtractDataFromGRPC(ctx)
			assert.True(t, ok)
			assert.Equal(t, v.expect, gwData)
		})
	}
}

func TestGatewayData_CheckAll(t *testing.T) {
	err := GatewayData{
		ProjectId:       "1",
		ProjectName:     "2",
		AccountIdentity: "3",
		AccountName:     "4",
	}.CheckAll()
	assert.NoError(t, err)

	err = GatewayData{}.CheckAll()
	assert.Contains(t, err.Error(), "ProjectId")
	assert.Contains(t, err.Error(), "ProjectName")
	assert.Contains(t, err.Error(), "AccountIdentity")
	assert.Contains(t, err.Error(), "AccountName")
}
