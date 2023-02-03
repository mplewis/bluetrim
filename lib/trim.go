package lib

import (
	"errors"
	"fmt"
	"os"
)

func Trim(video Metadata, out string, start float64, end float64) error {
	if _, err := os.Stat(out); err == nil {
		return errors.New("output file already exists")
	}
	// ffmpeg -i input.mp4 -ss 00:05:10 -to 00:15:30 -c:v copy -c:a copy output2.mp4
	out, _, err := call(
		"ffmpeg", "-hide_banner", "-loglevel", "error",
		"-i", video.Path,
		"-ss", fmt.Sprint(start), "-to", fmt.Sprint(end),
		"-c:v", "copy", "-c:a", "copy",
		out,
	)
	if err != nil {
		return fmt.Errorf("ffmpeg failed\noutput: %s\nerror: %w", out, err)
	}
	return nil
}
