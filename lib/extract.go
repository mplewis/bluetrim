package lib

import (
	"fmt"
	"os"
	"path"

	"github.com/sourcegraph/conc/iter"
	"github.com/thanhpk/randstr"
)

// Frame represents a single frame extracted from a video.
type Frame struct {
	Metadata
	Path string
}

// ExtractFrame extracts the frame at the given timestamp in seconds from the given video file into the target path.
func ExtractFrame(video Metadata, time float64, dest string) (Frame, error) {
	out, _, err := call("ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", fmt.Sprint(time), "-i", video.Path, "-frames:v", "1", dest)
	if err != nil {
		fmt.Println(out)
		return Frame{}, err
	}
	return Frame{
		Metadata: video,
		Path:     dest,
	}, nil
}

// extractFrames extracts the frames at the given timestamps from the given video file using ffmpeg.
func ExtractFrames(video Metadata, timestamps []float64) (string, []Frame, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", nil, err
	}
	frames, err := iter.MapErr[float64, Frame](timestamps, func(timestamp *float64) (Frame, error) {
		dest := path.Join(dir, fmt.Sprintf("frame_%f_%s.jpg", *timestamp, randstr.String(16)))
		if err != nil {
			return Frame{}, err
		}
		return ExtractFrame(video, *timestamp, dest)
	})
	return dir, frames, err
}
