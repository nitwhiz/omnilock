package server

import (
	"context"
	"github.com/nitwhiz/omnilock/pkg/table"
	"time"
)

type LockTable struct {
	locks *table.Table[string, uint64]
	ctx   context.Context
}

func NewLockTable(ctx context.Context) *LockTable {
	return &LockTable{
		locks: table.New[string, uint64](),
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
	return t.locks.TryPut(name, c.ID)
}

func (t *LockTable) Unlock(c *Client, name string) bool {
	return t.locks.RemoveIf(name, func(v uint64) bool {
		return c.ID == v
	})
}

func (t *LockTable) forceUnlock(name string) {
	t.locks.Remove(name)
}

func (t *LockTable) UnlockAllForClient(c *Client) {
	t.locks.RemoveByValue(c.ID)
}

func (t *LockTable) Count() int {
	return t.locks.Len()
}
