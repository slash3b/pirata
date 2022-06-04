package cache_test

import (
	"testing"

	"imdb/cache"

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

	c.Set("foo", 1)
	v, err = c.Get("foo")
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

func BenchmarkReplacementTest(b *testing.B) {
	c := cache.NewLRU[int, int](20) // what would be an optimal value ?
	for i := 0; i < b.N; i++ {
		c.Set(i, i)
	}
}
