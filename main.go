package main

import (
	"github.com/samchelini/swargo/bar"
  "log"
  "os"
)

func main() {
  log.SetOutput(os.Stderr) // log to stderr

  // initialize bar
	b := bar.NewBar()

  // create and add a DateTimeBlock
	dateTimeBlock := new(bar.DateTimeBlock)
	b.AddBlock(dateTimeBlock)

  // create and add a BrightnessBlock
  brightnessBlock := new(bar.BrightnessBlock)
  brightnessBlock.SetDir("/sys/class/backlight/intel_backlight")
  b.AddBlock(brightnessBlock)

  // run the bar
	b.Run()
}
