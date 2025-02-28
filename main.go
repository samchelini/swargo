package main

import (
	"github.com/samchelini/swargo/bar"
)

func main() {
	b := bar.NewBar()
  dateTimeBlock := new(bar.DateTimeBlock)
  b.AddBlock(dateTimeBlock)
  b.Run()
}
