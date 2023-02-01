package main

import (
	"fmt"
	"os"

	"github.com/mplewis/bluetrim/lib"
)

const similarThreshold = 0.03

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

	dir, frames, keyframe, err := lib.ExtractIntervalFrames(cfg, video, cfg.Interval.Seconds())
	if dir != "" {
		defer os.RemoveAll(dir)
	}
	check(err)

	similar, err := lib.SimilarFrames(similarThreshold, keyframe, frames)
	check(err)
	for _, frame := range similar {
		fmt.Printf("%s @ %f: %f\n", frame.Path, frame.Time, frame.Similarity)
	}

	// TODO: find blocks of contiguous dead air
	// TODO: binary search to find boundaries
}
