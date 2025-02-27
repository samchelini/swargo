package bar

import (
	"time"
)

type DateBlock struct {
	Block
}

func (b *DateBlock) Run() {
	for {
		t := time.Now().Format("01/02/2006 03:04:05 PM")
		b.FullText = t
		b.Update()
		time.Sleep(1 * time.Second)
	}
}

func (b *DateBlock) New() (*DateBlock) {
  return new(DateBlock)
}
