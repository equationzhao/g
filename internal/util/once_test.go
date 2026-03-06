package util

import (
	"sync"
	"testing"
)

func TestOnce_Do(t *testing.T) {
	var o Once
	var count int

	wg := sync.WaitGroup{}
	times := 10
	wg.Add(times)
	for range times {
		go func() {
			defer wg.Done()
			err := o.Do(func() error {
				count++
				return nil
			})
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		}()
	}
	wg.Wait()
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}
