package lib

import (
	"errors"
	"fmt"
	"os"
	"path"
)

// SearchForward searches forward from the given position at the given interval
// for the position of a frame that is not similar to the given keyframe.
func SearchForward(video Metadata, threshold float64, keyframe Frame, intervalSecs float64, pos float64) (float64, error) {
	last := pos
	fmt.Println(pos)
	for pos < video.DurationSeconds {
		pos += intervalSecs
		if pos > video.DurationSeconds {
			return last, nil
		}
		fmt.Println(pos)
		similar, err := similar(video, threshold, keyframe, pos)
		if err != nil {
			return 0, err
		}
		if !similar {
			return last, nil
		}
	}
	return 0, nil
}

// SearchBackward searches backward from the given position at the given interval
// for the position of a frame that is not similar to the given keyframe.
func SearchBackward(video Metadata, threshold float64, keyframe Frame, intervalSecs float64, pos float64) (float64, error) {
	last := pos
	for pos > 0 {
		pos -= intervalSecs
		if pos < 0 {
			return last, nil
		}
		fmt.Println(pos)
		similar, err := similar(video, threshold, keyframe, pos)
		if err != nil {
			return 0, err
		}
		if !similar {
			return last, nil
		}
	}
	return 0, nil
}

// Search searches for the position of a frame that is not similar to the given keyframe.
func Search(video Metadata, threshold float64, keyframe Frame, intervalSecs float64, pos float64) (float64, error) {
	if intervalSecs > 0 {
		return SearchForward(video, threshold, keyframe, intervalSecs, pos)
	}
	if intervalSecs < 0 {
		return SearchBackward(video, threshold, keyframe, -intervalSecs, pos)
	}
	return 0, errors.New("search interval must not be 0")
}

// similar returns true if the frame at the given position is similar to the given keyframe.
func similar(video Metadata, threshold float64, keyframe Frame, pos float64) (bool, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(dir)
	frame, err := ExtractFrame(video, pos, path.Join(dir, "frame.jpg"))
	if err != nil {
		return false, err
	}
	sim, err := CmpImages(keyframe.Path, frame.Path)
	if err != nil {
		return false, err
	}
	return sim < threshold, nil
}
