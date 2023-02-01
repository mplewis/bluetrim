package lib

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/sourcegraph/conc/iter"
)

type SimilarFrame struct {
	Frame
	Similarity float64
}

// magickCmpMatcher extracts the diff value from the output of Imagemagick's compare command.
var magickCmpMatcher = regexp.MustCompile(`^(\d+(\.\d+)?) \((.*)\)$`)

// cmpImages compares two images using Imagemagick's RMSE algorithm and returns a value from 0 to 1.
// A lower number indicates a higher similarity.
func CmpImages(a string, b string) (float64, error) {
	if a == b {
		return 0, nil
	}
	out, code, err := call("magick", "compare", "-metric", "RMSE", a, b, "NULL:")
	if err != nil && code != 1 {
		return 0, fmt.Errorf("magick compare failed\noutput: %s\nerror: %w", out, err)
	}
	matches := magickCmpMatcher.FindStringSubmatch(out)
	if matches == nil {
		return 0, fmt.Errorf("could not parse magick compare output: %s", out)
	}
	return strconv.ParseFloat(matches[3], 64)
}

func SimilarFrames(threshold float64, keyframe Frame, frames []Frame) ([]SimilarFrame, error) {
	sims, err := iter.MapErr(frames, func(frame *Frame) (float64, error) {
		return CmpImages(keyframe.Path, frame.Path)
	})
	if err != nil {
		return nil, err
	}
	var similar []SimilarFrame
	for i, sim := range sims {
		if sim < threshold {
			similar = append(similar, SimilarFrame{Frame: frames[i], Similarity: sim})
		}
	}
	return similar, nil
}
