package async

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewCoordinator(t *testing.T) {
	maxLoop := 100
	var wg sync.WaitGroup
	wg.Add(maxLoop)
	defer wg.Wait()

	consumer := func(ls []Element) error {
		fmt.Printf("get %+v \n", ls)
		time.Sleep(1 * time.Second)
		wg.Add(-len(ls))
		return nil
	}

	cfg := NewConfig(consumer,
		SetBatchSize(10),
		SetNumConsumer(2),
		SetBufferSize(maxLoop),
		SetBatchInterval(time.Second),
		SetRejectPolicy(Block),
	)
	c := NewCoordinator(cfg)
	c.Start()

	for i := 0; i < maxLoop; i++ {
		fmt.Printf("try put %v\n", i)
		discarded, err := c.Put(context.TODO(), i)
		if err != nil {
			fmt.Printf("discarded elements %+v for err %v", discarded, err)
			wg.Add(-len(discarded))
		}
	}
	c.Close(true)
}
