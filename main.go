package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

var (
	signalCounter int64
	prevSignals   int64
)

func main() {
	ch := make(chan os.Signal, 1000000)
	signal.Notify(ch)

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			nSignals := atomic.LoadInt64(&signalCounter)
			fmt.Printf("\r %d signals/s (%d total)", (nSignals-prevSignals)*4, nSignals)
			prevSignals = nSignals
		case sig := <-ch:
			switch sig {
			case os.Interrupt:
				return
			default:
				atomic.AddInt64(&signalCounter, 1)
			}
		}
	}
}
