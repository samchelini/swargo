package bar

import (
  "encoding/json"
  "github.com/samchelini/swargo/blocks"
)

type bar struct {
	header *header
	blocks []*blocks.Block
}

type header struct {
	Version     int  `json:"version"`
	ClickEvents bool `json:"click_events,omitempty"`
	ContSignal  int  `json:"cont_signal,omitempty"`
	StopSignal  int  `json:"stop_signal,omitempty"`
}

func NewBar() *bar {
	b := new(bar)
	b.header = &header{Version: 1}
  b.blocks = make([]*blocks.Block, 0)
	return b
}

func (b *bar) EnableClickEvents() {
	b.header.ClickEvents = true
}

func (b *bar) SetContSignal(signal int) {
	b.header.ContSignal = signal
}

func (b *bar) SetStopSignal(signal int) {
	b.header.StopSignal = signal
}

func (b *bar) String() string {
	json, _ := json.Marshal(b.header)
	return string(json)
}
