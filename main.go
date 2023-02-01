package main

import (
	"fmt"

	"github.com/mplewis/bluetrim/lib"
	"golang.org/x/exp/slices"
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

	half := video.DurationSeconds / 2
	pos := []float64{}
	for i := float64(0); i < half; i += cfg.Interval.Seconds() {
		pos = append(pos, float64(i))
	}
	for i := video.DurationSeconds; i > half; i -= cfg.Interval.Seconds() {
		pos = append(pos, float64(i))
	}
	slices.Sort(pos)

	dir, frames, err := lib.ExtractFrames(video, pos)
	// if dir != "" {
	// 	defer os.RemoveAll(dir)
	// }
	fmt.Println(dir)
	fmt.Println(frames)
	check(err)
}
