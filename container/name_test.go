package container

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostName(t *testing.T) {
	err := os.Setenv("HOSTNAME", "abcdefg")
	assert.NoError(t, err)

	assert.Equal(t, GetHostName(), "abcdefg")
}

func TestGetId(t *testing.T) {
	err := os.Setenv("HOSTNAME", "a-ab-ac-ad-ae-af-ag")
	assert.NoError(t, err)
	assert.Equal(t, GetId(), "ag")

	err = os.Setenv("HOSTNAME", "aabacadaeafag")
	assert.NoError(t, err)
	assert.Equal(t, GetId(), "aabacadaeafag")
}
