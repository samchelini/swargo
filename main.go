package main

import (
	"fmt"
	"github.com/samchelini/swargo/bar"
)

func main() {
	b := bar.NewBar()
	fmt.Println(b)

  dateBlock := new(bar.DateBlock)
  b.AddBlock(dateBlock)
  go dateBlock.Run()

  b.Run()
}
