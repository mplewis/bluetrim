package lib

import (
	"fmt"
	"os"
	"path"

	"github.com/schollz/progressbar/v3"
	"github.com/sourcegraph/conc/iter"
	"github.com/thanhpk/randstr"
	"golang.org/x/exp/slices"
)

// Frame represents a single frame extracted from a video.
type Frame struct {
	Metadata
	Path string
}

// ExtractFrame extracts the frame at the given timestamp in seconds from the given video file into the target path.
func ExtractFrameSecs(video Metadata, time float64, dest string) (Frame, error) {
	return ExtractFrameTimestamp(video, fmt.Sprint(time), dest)
}

// ExtractFrameTimestamp extracts the frame at the given FFmpeg timestamp from the given video file into the target path.
func ExtractFrameTimestamp(video Metadata, time string, dest string) (Frame, error) {
	out, _, err := call("ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", time, "-i", video.Path, "-frames:v", "1", dest)
	if err != nil {
		return Frame{}, fmt.Errorf("ffmpeg failed\noutput: %s\nerror: %w", out, err)
	}
	return Frame{
		Metadata: video,
		Path:     dest,
	}, nil
}

// ExtractFrames extracts the frames at the given timestamps from the given video file using ffmpeg.
func ExtractFrames(video Metadata, timestamps []float64) (string, []Frame, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", nil, err
	}
	pb := progressbar.Default(int64(len(timestamps)))
	frames, err := iter.MapErr[float64, Frame](timestamps, func(timestamp *float64) (Frame, error) {
		dest := path.Join(dir, fmt.Sprintf("frame_%f_%s.jpg", *timestamp, randstr.String(16)))
		if err != nil {
			return Frame{}, err
		}
		frame, err := ExtractFrameSecs(video, *timestamp, dest)
		pb.Add(1)
		return frame, err
	})
	return dir, frames, err
}

// ExtractIntervalFrames extracts frames from the given video file at the interval specified in the config.
// It also returns the keyframe.
func ExtractIntervalFrames(cfg Config, video Metadata, interval float64) (string, []Frame, Frame, error) {
	half := video.DurationSeconds / 2
	pos := []float64{}
	for i := float64(0); i < half; i += cfg.Interval.Seconds() {
		pos = append(pos, float64(i))
	}
	for i := video.DurationSeconds; i > half; i -= cfg.Interval.Seconds() {
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
		keyframe, err = ExtractFrameTimestamp(video, cfg.Keyframe, path.Join(dir, "keyframe.jpg"))
	}

	return dir, frames, keyframe, err
}
