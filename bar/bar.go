package bar

import (
	"encoding/json"
	"fmt"
)

// contains the swaybar header, array of blocks, and an update channel
type bar struct {
	header      *header
	blocks      []Block
	update      chan bool
	prettyPrint bool
}

// header defined by the swaybar protocol
type header struct {
	Version     int  `json:"version"`
	ClickEvents bool `json:"click_events,omitempty"`
	ContSignal  int  `json:"cont_signal,omitempty"`
	StopSignal  int  `json:"stop_signal,omitempty"`
}

// initializes and returns a new bar instance
func NewBar() *bar {
	b := new(bar)
	b.header = &header{Version: 1}
	b.blocks = make([]Block, 0)
	b.update = make(chan bool)
	return b
}

// set click events to true in the header
func (b *bar) EnableClickEvents() {
	b.header.ClickEvents = true
}

// enable pretty print for header and blocks JSON
func (b *bar) EnablePrettyPrint() {
	b.prettyPrint = true
}

// set continue signal in the header
func (b *bar) SetContSignal(signal int) {
	b.header.ContSignal = signal
}

// set stop signal in the header
func (b *bar) SetStopSignal(signal int) {
	b.header.StopSignal = signal
}

// returns header as a JSON string
func (b *bar) Header() string {
	var headerJSON []byte
	if b.prettyPrint {
		headerJSON, _ = json.MarshalIndent(b.header, "", "  ")
	} else {
		headerJSON, _ = json.Marshal(b.header)
	}
	return string(headerJSON)
}

// returns blocks as a JSON string
func (b *bar) Blocks() string {
	var blocksJSON []byte
	if b.prettyPrint {
		blocksJSON, _ = json.MarshalIndent(b.blocks, "  ", "  ")
	} else {
		blocksJSON, _ = json.Marshal(b.blocks)
	}
	return string(blocksJSON)
}

// sets update channel for the block and adds it to the bar
func (b *bar) AddBlock(block Block) {
	block.Sync(b.update)
	b.blocks = append(b.blocks, block)
}

// runs each block in a goroutine
func (b *bar) runBlocks() {
	for i := range b.blocks {
		go b.blocks[i].Run()
	}
}

// start running the blocks and wait for updates
func (b *bar) Run() {
	// start running all the blocks
	b.runBlocks()

	// print header and start of the infinite array
	fmt.Printf("%s%c[", b.Header(), 0x0A)
	if b.prettyPrint {
		fmt.Print("\n  ")
	}

	// print first array of blocks
	fmt.Print(b.Blocks())

	// continue printing blocks as they are updated in an infinite array
	for {
		<-b.update // wait for signal from the update channel

		// print blocks in infinite array
		fmt.Print(",")
		if b.prettyPrint {
			fmt.Print("\n  ")
		}
		fmt.Print(b.Blocks())
	}
}
