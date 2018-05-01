package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

type sigHistory struct {
	data []int64
}

func (s *sigHistory) Append(x int64) {
	if len(s.data) == cap(s.data) {
		copy(s.data, s.data[1:])
		s.data[len(s.data)-1] = x
	} else {
		s.data = append(s.data, x)
	}
}

func (s *sigHistory) Average() int64 {
	var total int64
	for _, v := range s.data {
		total += v
	}
	return total / int64(len(s.data))
}

var (
	history    = sigHistory{data: make([]int64, 0, 4)}
	sigCount   int64
	totalCount int64
)

func main() {
	ch := make(chan os.Signal, 1000000)
	signal.Notify(ch)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			nSignals := sigCount
			sigCount = 0
			totalCount += nSignals
			history.Append(nSignals)
			perSecond := history.Average()
			fmt.Printf("\r                                                                         ")
			fmt.Printf("\r%d signals/s\t\t\t%d total signals", perSecond, totalCount)
		case sig := <-ch:
			switch sig {
			case os.Interrupt:
				return
			default:
				sigCount++
			}
		}
	}
}
