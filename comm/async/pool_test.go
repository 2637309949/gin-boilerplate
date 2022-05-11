package async

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestWaitPoolExecutor_Wait(t *testing.T) {
	tr := NewWaitPool(2)
	var n int32 = 0
	for i := 0; i < 10; i++ {
		n := atomic.AddInt32(&n, 1)
		tr.Submit(func() {
			fmt.Printf("n=%d", n)
		})
	}
	tr.Wait()
}
