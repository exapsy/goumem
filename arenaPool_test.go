package goumem

import (
	"github.com/stretchr/testify/assert" // Assuming use of testify for assertions
	"testing"
)

func TestNewArenaPool(t *testing.T) {
	opts := ArenaPoolOptions{
		NumArenas: 5,
		ArenaSize: 1000, // Example size
	}

	ap, err := NewArenaPool(opts)
	assert.NoError(t, err)
	assert.NotNil(t, ap)
	assert.Equal(t, 5, len(ap.arenas))
}

func TestGetAndReturnArena(t *testing.T) {
	opts := ArenaPoolOptions{
		NumArenas: 2,
		ArenaSize: 1000,
	}

	ap, err := NewArenaPool(opts)
	assert.NoError(t, err)

	arena := ap.Get()
	assert.NotNil(t, arena)
	assert.True(t, arena.inUse)

	alloc, err := arena.Alloc(4)
	assert.NoError(t, err)
	assert.NotNil(t, alloc)

	err = arena.Free(alloc, 4)
	assert.NoError(t, err)

	ap.ReturnArena(arena)
	assert.False(t, arena.inUse)
}
