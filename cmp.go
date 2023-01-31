package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// magickCmpMatcher extracts the diff value from the output of Imagemagick's compare command.
var magickCmpMatcher = regexp.MustCompile(`^(\d+\.\d+) \((.*)\)$`)

// cmpImages compares two images using Imagemagick's RMSE algorithm and returns a value from 0 to 1.
// A lower number indicates a higher similarity.
func cmpImages(a string, b string) (float64, error) {
	out, code, err := call("magick", "compare", "-metric", "RMSE", a, b, "NULL:")
	if err != nil && code != 1 {
		return 0, err
	}
	fmt.Println(out)
	matches := magickCmpMatcher.FindStringSubmatch(out)
	if matches == nil {
		return 0, fmt.Errorf("could not parse magick compare output: %s", out)
	}
	return strconv.ParseFloat(matches[2], 64)
}

// extractFrames extracts the frames at the given timestamps from the given video file using ffmpeg.
func extractFrames(video string, timestamps []string) ([]Frame, error) {
	// TODO
	return nil, nil
}
