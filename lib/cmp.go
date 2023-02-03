package lib

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/sourcegraph/conc/iter"
)

// SimilarFrame is a frame with a similarity flag.
type SimilarFrame struct {
	Frame
	SimilarResult
}

type SimilarResult struct {
	Similar    bool
	Similarity float64
}

func (s SimilarFrame) String() string {
	return fmt.Sprintf("%f: %t (%f)", s.Time, s.Similar, s.Similarity)
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

// CmpFrames compares two frames using Imagemagick's RMSE algorithm and returns a value from 0 to 1.
// A lower number indicates a higher similarity.
func CmpFrames(a Frame, b Frame) (float64, error) {
	return CmpImages(a.Path, b.Path)
}

// CmpAllFrames compares all frames against the keyframe and returns whether they were similar to the keyframe.
func CmpAllFrames(keyframe Frame, threshold float64, frames []Frame) ([]SimilarFrame, error) {
	pb := progressbar.Default(int64(len(frames)))
	sims, err := iter.MapErr(frames, func(frame *Frame) (SimilarResult, error) {
		sim, err := CmpFrames(keyframe, *frame)
		pb.Add(1)
		return SimilarResult{Similarity: sim, Similar: sim < threshold}, err
	})
	if err != nil {
		return nil, err
	}
	sfs := make([]SimilarFrame, len(frames))
	for i, sim := range sims {
		sfs[i] = SimilarFrame{Frame: frames[i], SimilarResult: sim}
	}
	return sfs, nil
}
