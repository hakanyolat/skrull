package skrull

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMiddlewareStore_Get(t *testing.T) {
	store := &MiddlewareStore{data: map[string]interface{}{
		"x": 1,
	}}
	assert.Equal(t, 1, store.Get("x"))
}

func TestMiddlewareStore_Set(t *testing.T) {
	store := &MiddlewareStore{data: map[string]interface{}{}}
	store.Set("y", 2)
	assert.Equal(t, 2, store.Get("y"))

}

func TestMiddlewareStore_Has(t *testing.T) {
	store := &MiddlewareStore{data: map[string]interface{}{}}
	assert.False(t, store.Has("z"))

	store.Set("z", 3)
	assert.True(t, store.Has("z"))
}
