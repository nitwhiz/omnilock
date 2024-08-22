package lock

import (
	"context"
	"github.com/nitwhiz/omnilock/v2/pkg/util"
	"sync"
	"time"
)

type Table struct {
	locks *util.Map[string, uint64]
}

func NewTable() *Table {
	return &Table{
		locks: util.NewMap[string, uint64](),
	}
}

func (t *Table) TryLock(name string, clientId uint64) bool {
	return t.locks.TryPut(name, clientId)
}

func (t *Table) Lock(ctx context.Context, name string, clientId uint64) bool {
	if t.TryLock(name, clientId) {
		return true
	}

	result := false

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Microsecond * 250):
				if t.TryLock(name, clientId) {
					result = true
					return
				}

				break
			}
		}
	}()

	wg.Wait()

	return result
}

func (t *Table) Unlock(name string, clientId uint64) bool {
	return t.locks.RemoveIf(name, func(v uint64) bool {
		return clientId == v
	})
}

func (t *Table) UnlockAll(clientId uint64) {
	t.locks.RemoveByValue(clientId)
}
