package strsconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamel2Message(t *testing.T) {
	testTable := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "test1",
			in:   "TestName",
			out:  "Test name",
		},
		{
			name: "test2",
			in:   "testName",
			out:  "test name",
		},
		{
			name: "test3",
			in:   "testName3",
			out:  "test name3",
		},
		{
			name: "test4",
			in:   "test4Name",
			out:  "test4 name",
		},
		{
			name: "test5",
			in:   "TestName在",
			out:  "Test name在",
		},
	}

	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			assert.Equal(t, v.out, Camel2Message(v.in))
		})
	}
}
