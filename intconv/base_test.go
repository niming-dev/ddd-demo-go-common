package intconv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToBase(t *testing.T) {
	bs := []byte{255, 200, 155, 100, 55, 11, 0, 255, 200}
	str := BytesToBase(bs, Base62)
	assert.Equal(t, str, "LxWs8cX6Xt93E")

	str = BytesToBase(bs, Base32)
	assert.Equal(t, str, "x91r5co5osy89ge")
}
