package goumem

import "sync"

type ArenaPool struct {
	pool  *Pool
	inUse bool
	index int
	mutex sync.Mutex
}

func newArenaPool(arena *Pool, index int) *ArenaPool {
	return &ArenaPool{
		pool:  arena,
		index: index,
	}
}

func (ap *ArenaPool) use() {
	ap.inUse = true
}

func (ap *ArenaPool) free() {
	ap.inUse = false
}

func (ap *ArenaPool) Alloc(size uint) (*Ptr, error) {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	return ap.pool.Alloc(size)
}

func (ap *ArenaPool) Free(ptr *Ptr, size uint) error {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	return ap.pool.Free(ptr, size)
}

func newArenaPools(arenas []*Pool) []*ArenaPool {
	arenaPools := make([]*ArenaPool, 0)
	for i, arena := range arenas {
		arenaPools = append(arenaPools, newArenaPool(arena, i))
	}

	return arenaPools
}

type ArenaPools struct {
	arenas []*ArenaPool
	mutex  sync.Mutex
	next   int
}

type ArenaPoolOptions struct {
	NumArenas uint
	ArenaSize uint
}

func NewArenaPool(opts ArenaPoolOptions) (*ArenaPools, error) {
	arenas := make([]*Pool, 0)
	for i := uint(0); i < opts.NumArenas; i++ {
		pool, err := NewPool(PoolOptions{
			Size: opts.ArenaSize,
		})
		if err != nil {
			return nil, err
		}

		arenas = append(arenas, pool)
	}

	return &ArenaPools{
		arenas: newArenaPools(arenas),
	}, nil
}

func (ap *ArenaPools) Get() *ArenaPool {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	arena := ap.arenas[ap.next]
	arena.use()
	ap.next = (ap.next + 1) % len(ap.arenas)

	return arena
}

func (ap *ArenaPools) ReturnArena(arena *ArenaPool) {
	ap.mutex.Lock()
	defer ap.mutex.Unlock()

	arena.free()
	ap.next = arena.index
	return
}

func (ap *ArenaPools) Close() error {
	for _, arena := range ap.arenas {
		err := arena.pool.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
