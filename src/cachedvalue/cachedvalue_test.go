package cachedvalue_test

import (
	"fmt"
	"sync"
	"testing"

	"articles/src/cachedvalue"
)

func TestCachedValueMustBeConcurrencySafe(t *testing.T) {
	t.Parallel()

	type mode string

	var (
		modeEnable  mode = "enable"
		modeDisable mode = "disable"
		modeRandom  mode = "random"
	)

	type testcase struct {
		goroutineNum int
		mode         mode
	}

	testcases := []testcase{
		{1, modeEnable},
		{2, modeEnable},
		{3, modeEnable},
		{1, modeDisable},
		{2, modeDisable},
		{3, modeDisable},
		{1, modeRandom},
		{2, modeRandom},
		{3, modeRandom},
	}

	for _, tt := range testcases {
		name := fmt.Sprintf("gorourines:%d, mode:%s", tt.goroutineNum, tt.mode)

		t.Run(name, func(t *testing.T) {
			v := cachedvalue.NewCachedValue(42, func() int {
				return 100
			})

			var wg sync.WaitGroup
			defer wg.Wait()

			wg.Add(tt.goroutineNum)
			for goroutine := 0; goroutine < tt.goroutineNum; goroutine++ {
				goroutine := goroutine

				go func() {
					defer wg.Done()

					switch tt.mode {
					case modeEnable:
						v.Enable()
					case modeDisable:
						v.Disable()
					case modeRandom:
						if goroutine%2 == 0 {
							v.Enable()
						} else {
							v.Disable()
						}
					}
				}()
			}
		})
	}
}
