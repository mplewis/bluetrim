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
	fmt.Println(lib.Probe(cfg.In))
	video, err := lib.Probe(cfg.In)
	check(err)
	pos := []float64{}
	for i := 0; i < 120; i += 10 {
		pos = append(pos, float64(i))
	}
	dir, frames, err := lib.ExtractFrames(video, pos)
	// if dir != "" {
	// 	defer os.RemoveAll(dir)
	// }
	fmt.Println(dir)
	fmt.Println(frames)
	check(err)
}
