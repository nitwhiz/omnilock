package locking

import (
	"context"
	"github.com/nitwhiz/omnilock/pkg/client"
	"time"
)

type LockTable struct {
	locks *Table[string, uint64]
}

func NewLockTable() *LockTable {
	return &LockTable{
		locks: New[string, uint64](),
	}
}

func (t *LockTable) acquireLockWithContext(c *client.Client, name string, ctx context.Context) bool {
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

func (t *LockTable) Lock(c *client.Client, name string) bool {
	return t.acquireLockWithContext(c, name, c.GetContext())
}

func (t *LockTable) LockWithTimeout(c *client.Client, name string, timeout time.Duration) bool {
	lockCtx, cancel := context.WithTimeout(c.GetContext(), timeout)
	defer cancel()

	return t.acquireLockWithContext(c, name, lockCtx)
}

func (t *LockTable) TryLock(c *client.Client, name string) bool {
	return t.locks.TryPut(name, c.GetID())
}

func (t *LockTable) Unlock(c *client.Client, name string) bool {
	cID := c.GetID()

	return t.locks.RemoveIf(name, func(v uint64) bool {
		return cID == v
	})
}

func (t *LockTable) forceUnlock(name string) {
	t.locks.Remove(name)
}

func (t *LockTable) UnlockAllForClient(c *client.Client) {
	t.locks.RemoveByValue(c.GetID())
}

func (t *LockTable) Count() int {
	return t.locks.Len()
}
