package main

import (
	"fmt"
)

// Frame represents a single frame extracted from a video.
type Frame struct {
	Path  string
	Index int64
}

// cmpImages compares two images using Imagemagick's RMSE algorithm and returns a value.
// A lower number indicates a higher similarity.
func cmpImages(a string, b string) (float64, error) {
	// TODO
	return 0, nil
}

// extractFrames extracts the frames at the given timestamps from the given video file using ffmpeg.
func extractFrames(video string, timestamps []string) ([]Frame, error) {
	// TODO
	return nil, nil
}

// main runs the program.
func main() {
	cfg := LoadConfig()
	fmt.Println(probe(cfg.In))
}
