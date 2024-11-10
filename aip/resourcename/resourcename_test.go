package resourcename

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan_Struct(t *testing.T) {
	type dstStruct struct {
		Project   string `resourcename:"project"`
		Ticket    string `resourcename:"ticket"`
		Execution string `resourcename:"execution"`
		Node      string `resourcename:"node"`
	}

	table := []struct {
		tname   string
		name    string
		pattern string
		dst     dstStruct
		err     error
	}{
		{
			tname:   "normal1",
			name:    "projects/test_prj/tickets/test_ticket/executions/test_exe/nodes/test_node",
			pattern: "projects/{project}/tickets/{ticket}/executions/{execution}/nodes/{node}",
			dst: dstStruct{
				Project:   "test_prj",
				Ticket:    "test_ticket",
				Execution: "test_exe",
				Node:      "test_node",
			},
			err: nil,
		},
		{
			tname:   "normal2",
			name:    "projects/test_prj/tickets/test_ticket/executions/test_exe/setting",
			pattern: "projects/{project}/tickets/{ticket}/executions/{execution}/setting",
			dst: dstStruct{
				Project:   "test_prj",
				Ticket:    "test_ticket",
				Execution: "test_exe",
			},
			err: nil,
		},
		{
			tname:   "normal3",
			name:    "projects/test_prj/tickets/test_ticket/executions/test_exe/setting",
			pattern: "projects/{project}/tickets/{ticket}/executions/{execution}",
			dst:     dstStruct{},
			err:     errors.New("invalid pattern"),
		},
	}

	for _, v := range table {
		t.Run(v.tname, func(t *testing.T) {
			dst := dstStruct{}
			err := Scan(v.name, v.pattern, &dst)
			assert.Equal(t, err, v.err)
			assert.Equal(t, v.dst, dst)
		})
	}
}

func TestScan_MapSS(t *testing.T) {
	table := []struct {
		tname   string
		name    string
		pattern string
		dst     map[string]string
		err     error
	}{
		{
			tname:   "normal1",
			name:    "projects/test_prj/tickets/test_ticket/executions/test_exe/nodes/test_node",
			pattern: "projects/{project}/tickets/{ticket}/executions/{execution}/nodes/{node}",
			dst: map[string]string{
				"project":   "test_prj",
				"ticket":    "test_ticket",
				"execution": "test_exe",
				"node":      "test_node",
			},
			err: nil,
		},
		{
			tname:   "normal2",
			name:    "projects/test_prj/tickets/test_ticket/executions/test_exe/setting",
			pattern: "projects/{project}/tickets/{ticket}/executions/{execution}/setting",
			dst: map[string]string{
				"project":   "test_prj",
				"ticket":    "test_ticket",
				"execution": "test_exe",
			},
			err: nil,
		},
		{
			tname:   "normal3",
			name:    "projects/test_prj/tickets/test_ticket/executions/test_exe/setting",
			pattern: "projects/{project}/tickets/{ticket}/executions/{execution}",
			dst:     map[string]string{},
			err:     errors.New("invalid pattern"),
		},
	}

	for _, v := range table {
		t.Run(v.tname, func(t *testing.T) {
			dst := map[string]string{}
			err := Scan(v.name, v.pattern, &dst)
			assert.Equal(t, err, v.err)
			assert.Equal(t, v.dst, dst)
		})
	}
}
