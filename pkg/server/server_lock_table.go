package server

import (
	"context"
	"sync"
	"time"
)

type LockTable struct {
	mu    *sync.Mutex
	locks map[string]uint64
	ctx   context.Context
}

func NewLockTable(ctx context.Context) *LockTable {
	return &LockTable{
		mu:    &sync.Mutex{},
		locks: map[string]uint64{},
		ctx:   ctx,
	}
}

func (t *LockTable) acquireLockWithContext(c *Client, name string, ctx context.Context) bool {
	for {
		if t.TryLock(c, name) {
			return true
		}

		select {
		case <-ctx.Done():
			return false
		case <-time.After(time.Millisecond):
			break
		}
	}
}

func (t *LockTable) Lock(c *Client, name string) bool {
	return t.acquireLockWithContext(c, name, t.ctx)
}

func (t *LockTable) LockWithTimeout(c *Client, name string, timeout time.Duration) bool {
	lockCtx, cancel := context.WithTimeout(t.ctx, timeout)
	defer cancel()

	return t.acquireLockWithContext(c, name, lockCtx)
}

func (t *LockTable) TryLock(c *Client, name string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.locks[name]; ok {
		return false
	}

	t.locks[name] = c.ID

	return true
}

func (t *LockTable) Unlock(c *Client, name string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if cID, ok := t.locks[name]; ok && cID == c.ID {
		delete(t.locks, name)
		return true
	}

	return false
}

func (t *LockTable) forceUnlock(name string) {
	if _, ok := t.locks[name]; ok {
		delete(t.locks, name)
	}
}

func (t *LockTable) UnlockAllForClient(c *Client) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for lockName, cID := range t.locks {
		if cID == c.ID {
			t.forceUnlock(lockName)
		}
	}
}

func (t *LockTable) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return len(t.locks)
}
