package main

import (
	"fmt"
	"os"

	"github.com/mplewis/bluetrim/lib"
)

const similarThreshold = 0.20

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// main runs the program.
func main() {
	cfg := lib.LoadConfig()
	debug := func(x any) {
		if cfg.Debug {
			fmt.Println(x)
		}
	}

	fmt.Printf("Probing %s…\n", cfg.In)
	video, err := lib.Probe(cfg.In)
	check(err)

	fmt.Println("Extracting keyframe and summary frames…")
	dir, frames, keyframe, err := lib.ExtractFramesFull(cfg, video)
	defer os.RemoveAll(dir)
	check(err)
	sims, err := lib.CmpAllFrames(keyframe, similarThreshold, frames)
	check(err)
	for _, s := range sims {
		debug(s)
	}
	ranges := lib.WalkBoundaryRanges(sims)
	for _, r := range ranges {
		debug(r)
	}
	fmt.Printf("Found %d boundaries\n", len(ranges))

	frameInterval := 1.0
	fmt.Printf("Refining boundaries at %fs intervals…\n", frameInterval)
	var ranges2 []lib.TimeRangeState
	for _, r := range ranges {
		dir, frames, err := lib.ExtractIntervalFramesRange(video, frameInterval, r.Start, r.End)
		defer os.RemoveAll(dir)
		check(err)
		sims, err := lib.CmpAllFrames(keyframe, similarThreshold, frames)
		check(err)
		for _, s := range sims {
			debug(s)
		}
		ranges2 = append(ranges2, lib.WalkBoundaryRanges(sims)...)
	}
	for _, r := range ranges2 {
		debug(r)
	}

	frameInterval = 1 / video.FrameRate
	fmt.Printf("Refining boundaries at %fs (frame) intervals…\n", frameInterval)
	var ranges3 []lib.TimeRangeState
	for _, r := range ranges2 {
		dir, frames, err := lib.ExtractIntervalFramesRange(video, frameInterval, r.Start, r.End)
		defer os.RemoveAll(dir)
		check(err)
		sims, err := lib.CmpAllFrames(keyframe, similarThreshold, frames)
		check(err)
		for _, s := range sims {
			debug(s)
		}
		ranges3 = append(ranges3, lib.WalkBoundaryRanges(sims)...)
	}
	for _, r := range ranges3 {
		debug(r)
	}

	keepers := lib.TimeRangesToKeepers(ranges3)
	if len(keepers) == 0 {
		fmt.Println("No keepers, not trimming.")
		return
	}

	for i, r := range keepers {
		if r.End > video.DurationSeconds {
			keepers[i].End = video.DurationSeconds
		}
	}

	fmt.Println("Keeping the following time ranges:")
	for _, r := range keepers {
		fmt.Println(r)
	}

	if cfg.DryRun {
		fmt.Println("Dry run, not trimming.")
		return
	}

	if len(keepers) == 1 {
		fmt.Printf("Trimming video to %s…\n", cfg.Out)
		err = lib.Trim(video, cfg.Out, keepers[0].Start, keepers[0].End)
		check(err)
		return
	}

	for i, keeper := range keepers {
		out := lib.SuffixFn(cfg.Out, fmt.Sprintf("_%d", i+1))
		fmt.Printf("Trimming clip %d to %s…\n", i+1, out)
		err := lib.Trim(video, out, keeper.Start, keeper.End)
		check(err)
	}
}
