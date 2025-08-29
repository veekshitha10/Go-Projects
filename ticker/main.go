package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func main() {
	tickers := []string{"AAPL", "GOOG", "INFY"}
	ch1 := make(chan float64)
	ch2 := make(chan float64)
	ch3 := make(chan float64)
	done := make(chan struct{})

	go sender(ch1, done)
	go sender(ch2, done)
	go sender(ch3, done)

	RececeiverBySelectClose(ch1, ch2, ch3, tickers)

	close(done)
}

func sender(c chan float64, done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case <-time.After(time.Second):
			c <- rand.Float64()
		}
	}
}

func After(duration time.Duration) chan struct{} {
	aft := make(chan struct{})
	go func() {
		time.Sleep(duration)
		aft <- struct{}{}
		close(aft)
	}()
	return aft
}

func RececeiverBySelectClose(ch1, ch2, ch3 chan float64, tickers []string) {
	aft := After(10 * time.Second)

	for {
		select {
		case v := <-ch1:
			fmt.Printf("%s %s : %.2f\n", time.Now().Format("15:04:05"), tickers[0], v*1000)
		case v := <-ch2:
			fmt.Printf("%s %s : %.2f\n", time.Now().Format("15:04:05"), tickers[1], v*1000)
		case v := <-ch3:
			fmt.Printf("%s %s : %.2f\n", time.Now().Format("15:04:05"), tickers[2], v*1000)
		case <-aft:
			return
		}
	}
}
