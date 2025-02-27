package bar

import (
	"encoding/json"
  "fmt"
)

type bar struct {
	header  *header
	blocks  []BlockRunner
	updater chan bool
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
	b.blocks = make([]BlockRunner, 0)
	b.updater = make(chan bool)
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

func (b *bar) AddBlock(block BlockRunner) {
	block.Sync(b.updater)
	b.blocks = append(b.blocks, block)
}

func (b *bar) Run() {
  for {
    _ = <-b.updater
    for i := range b.blocks {
      fmt.Println(b.blocks[i])
    }
  }
}
