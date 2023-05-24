package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	//ticker, stop := newTickEveryMinute()
	//ticker, stop := NewTickEveryMinute(3)
	ticker := NewTickEveryMinute1([]uint8{3, 23, 37})

	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.Ch:
			fmt.Printf("tick at %v\n", time.Now().Second())
		}
	}
}

func newTickEveryMinute() (<-chan struct{}, func()) {
	ticker_3_second := time.NewTicker(3 * time.Second)
	ticker_10_second := time.NewTicker(10 * time.Second)
	ticker_1_minute := time.NewTicker(1 * time.Minute)

	ch := make(chan struct{}, 0)
	exitSignal := make(chan struct{}, 0)

	stop := func() {
		close(exitSignal)
		ticker_3_second.Stop()
		ticker_10_second.Stop()
		ticker_1_minute.Stop()
	}

	go func() {
		for {
			select {
			case <-exitSignal:
				return
			case <-ticker_3_second.C:
				ch <- struct{}{}
			case <-ticker_10_second.C:
				ch <- struct{}{}
			case <-ticker_1_minute.C:
				ch <- struct{}{}
			}
		}
	}()

	return ch, stop
}

func NewTickEveryMinute(second int) (chan<- struct{}, func()) {

	now := time.Now()
	sleepDur := time.Duration((60-now.Second()+second)%60) * time.Second

	// we sleep here for sleepDur before creating our ticker
	time.Sleep(sleepDur)

	// we create a tickers
	ticker_3_Second := time.NewTicker(1 * time.Minute)

	ch := make(chan<- struct{}, 0)
	exitSignal := make(chan struct{}, 0)

	// we call stop from parent to stop our ticker
	stop := func() {
		// closing exitSignal will return from below go routine
		close(exitSignal)
		ticker_3_Second.Stop()
	}

	go func() {
		for {
			select {
			case <-exitSignal:
				return
			case <-ticker_3_Second.C:
				ch <- struct{}{}
			}
		}
	}()

	return ch, stop
}

type CustomTicker struct {
	Ch         chan struct{}
	Tickers    []*time.Ticker
	ExitSignal chan struct{}
}

// Stop below will stop all tickers in tickers slice
func (c *CustomTicker) Stop() {
	c.ExitSignal <- struct{}{}
	for _, t := range c.Tickers {
		t.Stop()
	}
}

// Add we add a ticker to tickers slice
func (c *CustomTicker) Add(ticker *time.Ticker) {
	c.Tickers = append(c.Tickers, ticker)
}

func NewTickEveryMinute1(seconds []uint8) *CustomTicker {

	customTicker := CustomTicker{
		Ch: make(chan struct{}),
	}

	for _, second := range seconds {
		if second > 60 {
			continue
		}

		// for each second we create a new ticker in separate goroutines
		// and all the ticker pass on their ticks to our customTicker.Ch

		go func(second_ uint8, t *CustomTicker) {
			now := time.Now()
			sleepDur := time.Duration((60-now.Second()+int(second_))%60) * time.Second

			time.Sleep(sleepDur)

			ticker := time.NewTicker(time.Minute)
			t.Add(ticker)

			for {
				select {
				case <-t.ExitSignal:
					return
				case <-ticker.C:
					t.Ch <- struct{}{}
				}
			}
		}(second, &customTicker)

	}
	return &customTicker
}
