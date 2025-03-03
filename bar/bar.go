package bar

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	lflags int    = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix
	green  string = ": \033[92m"
	red    string = ": \033[91m"
	reset  string = "\033[0m"
)

// contains the swaybar header, array of blocks, and an update channel
type bar struct {
	header      *header
	blocks      []Block
	update      chan bool
	err         chan string
	prettyPrint bool
	logger      *log.Logger
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
	b.err = make(chan string)
	return b
}

// set click events to true in the header
func (b *bar) EnableClickEvents() {
	b.header.ClickEvents = true
}

// enable logging to stderr for debugging
func (b *bar) EnableLogging() {
	b.logger = log.New(os.Stderr, "", lflags)
	b.Log("logging enabled")
}

// log any message (green colored)
func (b *bar) Log(msg string) {
	if b.logger != nil {
		b.logger.SetPrefix(green)
		b.logger.Printf("%s%s\n", msg, reset)
	}
}

// log an error (red colored)
func (b *bar) LogError(msg string) {
	if b.logger != nil {
		b.logger.SetPrefix(red)
		b.logger.Printf("%s%s\n", msg, reset)
	}
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
	var err error
	var headerJSON []byte
	if b.prettyPrint {
		headerJSON, err = json.MarshalIndent(b.header, "", "  ")
	} else {
		headerJSON, err = json.Marshal(b.header)
	}
	if err != nil {
		b.LogError("error marshalling header: " + err.Error())
	}

	return string(headerJSON)
}

// returns blocks as a JSON string
func (b *bar) Blocks() string {
	var err error
	var blocksJSON []byte
	if b.prettyPrint {
		blocksJSON, err = json.MarshalIndent(b.blocks, "  ", "  ")
	} else {
		blocksJSON, err = json.Marshal(b.blocks)
	}
	if err != nil {
		b.LogError("error marshalling blocks: " + err.Error())
	}
	return string(blocksJSON)
}

// sets update and error channel for the block and adds it to the bar
func (b *bar) AddBlock(block Block) {
	block.Sync(b.update, b.err)
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
		select {
		case <-b.update: // wait for signal from the update channel

			// print blocks in infinite array
			fmt.Print(",")
			if b.prettyPrint {
				fmt.Print("\n  ")
			}
			fmt.Print(b.Blocks())
		case msg := <-b.err: // handle errors from blocks
			b.LogError(msg)
		}
	}
}
