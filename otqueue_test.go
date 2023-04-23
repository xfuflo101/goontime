package goontime_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	got "github.com/xfuflo101/goontime"
)

func TestOTQueue01(t *testing.T) {

	t.Logf("TestOTQueue01 - started")

	var wg sync.WaitGroup
	ctx, ctxCancel := context.WithCancel(context.Background())

	onTimeQueue := got.NewOnTimeQueue[string](10)
	defer onTimeQueue.Shutdown()

	onTimeChan := got.CreateOnTimeQueueChannel(ctx, &wg, onTimeQueue)
	// this channel will close writer by ctx cancel

	start := time.Now()

	var resultStr string

	wg.Add(1)
	go func() {
		// reader
		defer wg.Done()
		defer func() {
			t.Logf("TestUtils01 - now[%v] - reader is out of loop", time.Now())
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case str := <-onTimeChan:
				now := time.Now()
				dur := now.Sub(start)
				resultStr += fmt.Sprintf("; [%v : %v]", str, dur.Milliseconds()/100)
				t.Logf("TestUtils01 - now[%v] - str[%v]", now, str)
			}
		}
	}()

	onTimeQueue.Add("10", time.Duration(10*time.Second))
	onTimeQueue.Add("8", time.Duration(8*time.Second))
	onTimeQueue.Add("3", time.Duration(3*time.Second))
	onTimeQueue.Add("2", time.Duration(2*time.Second))
	onTimeQueue.Add("8_2", time.Duration(8*time.Second))
	onTimeQueue.Add("1", time.Duration(1*time.Second))

	time.Sleep(1 * time.Second)

	onTimeQueue.Add("12", time.Duration(12*time.Second))
	onTimeQueue.Add("7", time.Duration(7*time.Second))
	onTimeQueue.Add("15", time.Duration(15*time.Second))
	onTimeQueue.Add("8_3", time.Duration(8*time.Second))
	onTimeQueue.Add("3_2", time.Duration(3*time.Second))
	onTimeQueue.Add("2_2", time.Duration(2*time.Second))
	onTimeQueue.Add("1_2", time.Duration(1*time.Second))

	time.Sleep(17 * time.Second)

	ctxCancel()

	wg.Wait()

	assert.Equal(t,
		"; [1 : 10]; [2 : 20]; [1_2 : 20]; [3 : 30]; [2_2 : 30]; [3_2 : 40]; [8 : 80]; [8_2 : 80]; [7 : 80]; [8_3 : 90]; [10 : 100]; [12 : 130]; [15 : 160]",
		resultStr, "TestOTQueue01",
	)

	t.Logf("TestOTQueue01 - OK")
}
