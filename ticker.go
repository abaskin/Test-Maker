package testparts

import "time"

type Ticker struct {
	ticker      *time.Ticker
	stopChannel chan bool
	stepFunc    func()
	stopFunc    func()
}

func NewTicker(tickTime time.Duration, stop, step func()) *Ticker {
	tick := &Ticker{
		ticker:      time.NewTicker(tickTime),
		stopChannel: make(chan bool),
		stepFunc:    step,
		stopFunc:    stop,
	}
	go func() {
		for {
			select {
			case <-tick.stopChannel:
				if stop != nil {
					stop()
				}
				return
			case <-tick.ticker.C:
				if step != nil {
					step()
				}
			}
		}
	}()
	return tick
}

func (t *Ticker) Stop() {
	t.stopChannel <- true
}
