package lib

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Frame represents a single frame extracted from a video.
type Frame struct {
	Metadata
	Path  string
	Index int64
}

// tsParser is a parser for a timestamp format.
type tsParser struct {
	format *regexp.Regexp
	parser func([]string) (time.Duration, error)
}

// TargetTime represents a point in time in a video.
type TargetTime struct {
	// exactly one is required
	Timestamp  string
	FrameIndex int64
}

// tsParsers is a list of all supported parsers for known timestamp formats.
var tsParsers = []tsParser{
	{
		format: regexp.MustCompile(`^(\d\d?):(\d\d):(\d\d)(.(\d+))?$`),
		parser: func(matches []string) (time.Duration, error) {
			h, err := strconv.Atoi(matches[1])
			if err != nil {
				return 0, err
			}
			m, err := strconv.Atoi(matches[2])
			if err != nil {
				return 0, err
			}
			s, err := strconv.Atoi(matches[3])
			if err != nil {
				return 0, err
			}
			ms := 0
			if len(matches) > 4 {
				ms, err = strconv.Atoi(matches[5])
				if err != nil {
					return 0, err
				}
			}
			return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second + time.Duration(ms)*time.Millisecond, nil
		},
	},

	{
		format: regexp.MustCompile(`^(\d\d?):(\d\d)(.(\d+))?$`),
		parser: func(matches []string) (time.Duration, error) {
			m, err := strconv.Atoi(matches[1])
			if err != nil {
				return 0, err
			}
			s, err := strconv.Atoi(matches[2])
			if err != nil {
				return 0, err
			}
			ms := 0
			if len(matches) > 3 {
				ms, err = strconv.Atoi(matches[4])
				if err != nil {
					return 0, err
				}
			}
			return time.Duration(m)*time.Minute + time.Duration(s)*time.Second + time.Duration(ms)*time.Millisecond, nil
		},
	},

	{
		format: regexp.MustCompile(`^(\d\d?)(.(\d+))?$`),
		parser: func(matches []string) (time.Duration, error) {
			s, err := strconv.Atoi(matches[1])
			if err != nil {
				return 0, err
			}
			ms := 0
			if len(matches) > 2 {
				ms, err = strconv.Atoi(matches[3])
				if err != nil {
					return 0, err
				}
			}
			return time.Duration(s)*time.Second + time.Duration(ms)*time.Millisecond, nil
		},
	},
}

// ParseTs parses a timestamp string into a time.Duration.
func ParseTs(ts string) (time.Duration, error) {
	for _, p := range tsParsers {
		matches := p.format.FindStringSubmatch(ts)
		if len(matches) == 0 {
			continue
		}
		dur, err := p.parser(matches)
		if err != nil {
			return 0, err
		}
		return dur, nil
	}
	return 0, fmt.Errorf("could not parse timestamp: %s", ts)
}

// TsToFrame converts a timestamp	string to a frame number.
func TsToFrame(ts string, fps float64) (int64, error) {
	dur, err := ParseTs(ts)
	if err != nil {
		return 0, err
	}
	return int64(dur.Seconds()*fps) + 1, nil
}

func ExtractFrame(video Metadata, time TargetTime, dest string) (string, error) {
	ts := ts
	fnum, err := video.TsToFrame(ts)
	if err != nil {
		return nil, err
	}

	_, err := call("ffmpeg", "-hide_banner", "-loglevel", "error", "-ss", ts, "-i", video.Path, "-frames:v", "1", path)
	path := "TODO"
	frame := Frame{
		Metadata: video,
		Path:     path,
		Index:    fnum,
	}
}

// extractFrames extracts the frames at the given timestamps from the given video file using ffmpeg.
func ExtractFrames(video Metadata, timestamps []string) ([]Frame, error) {
	var frames []Frame
	for _, ts := range timestamps {
		ts := ts
		fnum, err := video.TsToFrame(ts)
		if err != nil {
			return nil, err
		}
		// TODO: extract frame
		path := "TODO"
		frames = append(frames, Frame{
			Metadata: video,
			Path:     path,
			Index:    fnum,
		})
	}
	return frames, nil
}
