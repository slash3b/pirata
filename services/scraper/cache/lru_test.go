package cache_test

import (
	"scraper/cache"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicReplacementTest(t *testing.T) {
	c := cache.NewLRU[string, int](3)
	res, err := c.Get("a")
	assert.Error(t, err)
	assert.Equal(t, 0, res)

	c.Set("foo", 1)
	v, err := c.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	c.Set("bar", 2)
	v, err = c.Get("bar")
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	c.Set("baz", 3)
	v, err = c.Get("baz")
	assert.NoError(t, err)
	assert.Equal(t, 3, v)

	c.Set("barf", 4)
	v, err = c.Get("barf")
	assert.NoError(t, err)
	assert.Equal(t, 4, v)

	v, err = c.Get("foo")
	assert.Error(t, err)
	assert.Equal(t, 0, v)
}
