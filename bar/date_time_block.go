package bar

import (
	"time"
)

// block for displaying date and time
type DateTimeBlock struct {
	BlockTemplate
}

// gets and formats the current time, updates every second
func (b *DateTimeBlock) Run() {
	for {
		t := time.Now().Format("01/02/2006 03:04:05 PM")
		b.FullText = t
		b.Update()
		time.Sleep(1 * time.Second)
	}
}

// returns a new instance of DateTimeBlock
func (b *DateTimeBlock) New() (*DateTimeBlock) {
  return new(DateTimeBlock)
}
