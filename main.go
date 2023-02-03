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
		fmt.Printf("%s @ %f: %t\n", frame.Path, frame.Time, frame.Similar)
	}
	trs := lib.PartitionTimeRanges(similar)
	for _, tr := range trs {
		fmt.Printf("%f - %f\n", tr.Start, tr.End)
	}

	// TODO: binary search to find boundaries
}
