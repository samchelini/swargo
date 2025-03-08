package bar

import (
	"time"
)

// block for displaying date and time
type DateTimeBlock struct {
	BlockTemplate
}

// updates the time every second
func (block *DateTimeBlock) Run() {
	for {
		block.SetFullText(block.getTime())
		block.Update()
		time.Sleep(1 * time.Second)
	}
}

// get formatted time as string
func (block *DateTimeBlock) getTime() string {
	return time.Now().Format("01/02/2006 03:04:05 PM")
}
