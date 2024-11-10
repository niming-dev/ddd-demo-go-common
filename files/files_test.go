package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllFile(t *testing.T) {
	fs, err := GetAllFile(".", ".go")
	assert.NoError(t, err)
	assert.Equal(t, []string{"./files.go", "./files_test.go"}, fs)

	fs, err = GetAllFile(".", ".g")
	assert.NoError(t, err)
	assert.Equal(t, []string(nil), fs)

	fs, err = GetAllFile(".")
	assert.NoError(t, err)
	assert.Equal(t, []string{"./files.go", "./files_test.go"}, fs)
}
