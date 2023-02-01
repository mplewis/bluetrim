package main

import (
	"fmt"

	"github.com/mplewis/bluetrim/lib"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// main runs the program.
func main() {
	cfg := lib.LoadConfig()
	video, err := lib.Probe(cfg.In)
	check(err)

	dir, frames, err := lib.ExtractIntervalFrames(cfg, video, cfg.Interval.Seconds())
	// if dir != "" {
	// 	defer os.RemoveAll(dir)
	// }
	fmt.Println(dir)
	fmt.Println(frames)
	check(err)
}
