package pagetoken

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePageToken_simply(t *testing.T) {
	pt := New(0, 100)

	pt, err := ParsePageToken(pt.String())
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), pt.Offset)
	assert.Equal(t, uint32(100), pt.Limit)
	assert.Equal(t, "", pt.OrderBy)
	assert.Equal(t, "", pt.Filter)
	assert.Equal(t, "", pt.Parent)
}

func TestParsePageToken_withParent(t *testing.T) {
	pt := New(0, 100)
	pt.Parent = "test_parent"

	pt, err := ParsePageToken(pt.String())
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), pt.Offset)
	assert.Equal(t, uint32(100), pt.Limit)
	assert.Equal(t, "", pt.OrderBy)
	assert.Equal(t, "", pt.Filter)
	assert.Equal(t, "test_parent", pt.Parent)
}

func TestParsePageToken_withAll(t *testing.T) {
	pt := New(0, 100)
	pt.Parent = "test_parent"
	pt.OrderBy = "id,create_at desc"
	pt.Filter = `name="abc"`

	pt, err := ParsePageToken(pt.String())
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), pt.Offset)
	assert.Equal(t, uint32(100), pt.Limit)
	assert.Equal(t, "id,create_at desc", pt.OrderBy)
	assert.Equal(t, `name="abc"`, pt.Filter)
	assert.Equal(t, "test_parent", pt.Parent)
}

func TestParsePageToken_withAll_maxLimit(t *testing.T) {
	pt := New(0, 100, WithMaxLimit(50))
	pt.Parent = "test_parent"
	pt.OrderBy = "id,create_at desc"
	pt.Filter = `name="abc"`

	pt, err := ParsePageToken(pt.String())
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), pt.Offset)
	assert.Equal(t, uint32(50), pt.Limit)
	assert.Equal(t, "id,create_at desc", pt.OrderBy)
	assert.Equal(t, `name="abc"`, pt.Filter)
	assert.Equal(t, "test_parent", pt.Parent)
}

func TestParsePageToken_withAll_defaultOrderBy(t *testing.T) {
	pt := New(0, 100, WithDefaultOrderBy("id,create_at desc"))
	pt.Parent = "test_parent"
	pt.Filter = `name="abc"`

	pt, err := ParsePageToken(pt.String())
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), pt.Offset)
	assert.Equal(t, uint32(100), pt.Limit)
	assert.Equal(t, "id,create_at desc", pt.OrderBy)
	assert.Equal(t, `name="abc"`, pt.Filter)
	assert.Equal(t, "test_parent", pt.Parent)
}
