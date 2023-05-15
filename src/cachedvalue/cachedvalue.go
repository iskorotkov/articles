package cachedvalue

import (
	"context"
	"sync/atomic"
	"time"
)

func NewCachedValue[T any](initialValue T, refresh func() T) *CachedValue[T] {
	v := &CachedValue[T]{
		value:   atomic.Value{},
		enabled: atomic.Bool{},
		refresh: refresh,
		cancel:  nil,
	}
	v.Enable()

	return v
}

// CachedValue is a wrapper around a value of type T that automatically refreshes it every second.
type CachedValue[T any] struct {
	value   atomic.Value
	enabled atomic.Bool
	refresh func() T
	cancel  context.CancelFunc
}

// Get returns latest value.
func (v *CachedValue[T]) Get() T {
	return v.value.Load().(T)
}

// Enable turns on automatic refresh.
//
// It is safe to use concurrently, but not with Disable.
func (v *CachedValue[T]) Enable() {
	if v.enabled.Swap(true) {
		return // Already enabled.
	}

	// Setup cancellation for background refresh goroutine.
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	go func() {
		timer := time.NewTimer(time.Second)
		for {
			select { // Context done or time to refresh, whichever happens first.
			case <-ctx.Done():
				return
			case <-timer.C:
				newValue := v.refresh()
				v.value.Store(newValue)
			}
		}
	}()
}

// Disable turns off automatic refresh.
//
// It is safe to use concurrently, but not with Enable.
func (v *CachedValue[T]) Disable() {
	if !v.enabled.Swap(false) {
		return // Already disabled.
	}

	if v.cancel == nil {
		panic("enabled CachedValue must have non-nil v.cancel")
	}

	v.cancel()
	v.cancel = nil // Remove reference to CancelFunc, that references Context, that may reference lots of things.
}
