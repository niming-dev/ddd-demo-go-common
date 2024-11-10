package shortuuid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	id := New()
	assert.True(t, len(id) == 22)

	fmt.Println("ShortUUID:", id)
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}
