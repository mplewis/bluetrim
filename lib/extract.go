package lib

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/sourcegraph/conc/iter"
	"github.com/thanhpk/randstr"
	"golang.org/x/exp/slices"
)

// Frame represents a single frame extracted from a video.
type Frame struct {
	Metadata
	Path string
	Time float64
}

// ExtractFrame extracts the frame at the given timestamp in seconds from the given video file into the target path.
func ExtractFrame(video Metadata, time float64, dest string) (Frame, error) {
	out, _, err := call("ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", fmt.Sprint(time), "-i", video.Path, "-frames:v", "1", dest)
	if err != nil {
		return Frame{}, fmt.Errorf("ffmpeg failed\noutput: %s\nerror: %w", out, err)
	}
	return Frame{
		Metadata: video,
		Path:     dest,
		Time:     time,
	}, nil
}

// ExtractFrames extracts the frames at the given timestamps from the given video file using ffmpeg.
func ExtractFrames(video Metadata, timestamps []float64) (string, []Frame, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", nil, err
	}
	pb := progressbar.Default(int64(len(timestamps)))
	frames, err := iter.MapErr(timestamps, func(timestamp *float64) (Frame, error) {
		dest := path.Join(dir, fmt.Sprintf("frame_%s.jpg", randstr.String(16)))
		if err != nil {
			return Frame{}, err
		}
		frame, err := ExtractFrame(video, *timestamp, dest)
		pb.Add(1)
		return frame, err
	})
	return dir, frames, err
}

// ExtractIntervalFrames extracts frames from the given video file at the interval specified in the config.
// It also returns the keyframe.
func ExtractIntervalFrames(cfg Config, video Metadata, interval float64) (string, []Frame, Frame, error) {
	var err error
	keyframeTs := float64(-1)
	if cfg.Keyframe != "start" && cfg.Keyframe != "end" {
		keyframeTs, err = strconv.ParseFloat(cfg.Keyframe, 64)
		if err != nil {
			return "", nil, Frame{}, fmt.Errorf("invalid keyframe %q: %w", cfg.Keyframe, err)
		}
	}

	half := video.DurationSeconds / 2
	pos := []float64{}
	for i := float64(0); i < half; i += cfg.Interval.Seconds() {
		pos = append(pos, float64(i))
	}
	for i := video.DurationSeconds - 1; i > half; i -= cfg.Interval.Seconds() {
		pos = append(pos, float64(i))
	}
	slices.Sort(pos)

	dir, frames, err := ExtractFrames(video, pos)
	if err != nil {
		return "", nil, Frame{}, err
	}

	var keyframe Frame
	if cfg.Keyframe == "start" {
		keyframe = frames[0]
	} else if cfg.Keyframe == "end" {
		keyframe = frames[len(frames)-1]
	} else {
		keyframe, err = ExtractFrame(video, keyframeTs, path.Join(dir, "keyframe.jpg"))
	}

	return dir, frames, keyframe, err
}
