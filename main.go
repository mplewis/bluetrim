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

	fmt.Println("Probing video…")
	video, err := lib.Probe(cfg.In)
	check(err)

	fmt.Println("Extracting keyframe and summary frames…")
	dir, frames, keyframe, err := lib.ExtractFramesFull(cfg, video)
	defer os.RemoveAll(dir)
	check(err)
	sims, err := lib.CmpAllFrames(keyframe, similarThreshold, frames)
	check(err)
	ranges := lib.WalkBoundaryRanges(sims)

	frameInterval := 1.0
	fmt.Printf("Searching for boundaries at %fs intervals…\n", frameInterval)
	var ranges2 []lib.TimeRangeState
	for _, r := range ranges {
		dir, frames, err := lib.ExtractIntervalFramesRange(video, frameInterval, r.Start, r.End)
		defer os.RemoveAll(dir)
		check(err)
		sims, err := lib.CmpAllFrames(keyframe, similarThreshold, frames)
		check(err)
		ranges2 = append(ranges2, lib.WalkBoundaryRanges(sims)...)
	}

	frameInterval = 1 / video.FrameRate
	fmt.Printf("Searching for boundaries at %fs (frame) intervals…\n", frameInterval)
	var ranges3 []lib.TimeRangeState
	for _, r := range ranges2 {
		dir, frames, err := lib.ExtractIntervalFramesRange(video, frameInterval, r.Start, r.End)
		defer os.RemoveAll(dir)
		check(err)
		sims, err := lib.CmpAllFrames(keyframe, similarThreshold, frames)
		check(err)
		ranges3 = append(ranges3, lib.WalkBoundaryRanges(sims)...)
	}

	keepers := lib.TimeRangesToKeepers(ranges3)
	fmt.Println(keepers)
}
