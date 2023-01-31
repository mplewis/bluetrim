package main

import (
	"fmt"
)

// Frame represents a single frame extracted from a video.
type Frame struct {
	Path  string
	Index int64
}

// main runs the program.
func main() {
	cfg := LoadConfig()
	fmt.Println(probe(cfg.In))
	fmt.Println(cmpImages("/Users/mplewis/tmp/framex/output_0001.jpg", "/Users/mplewis/tmp/framex/output_0002.jpg"))
}
