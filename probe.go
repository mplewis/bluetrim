package main

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// frameRateMatcher extracts the numerator and denominator from a frame rate.
var frameRateMatcher = regexp.MustCompile(`^(\d+)/(\d+)$`)

// Metadata represents the metadata for a video.
type Metadata struct {
	FrameCount      int64
	DurationSeconds float64
}

// ProbeOutputStream represents the metadata for a single stream in a video.
type ProbeOutputStream struct {
	CodecType string `json:"codec_type"`     // video, audio
	FrameRate string `json:"avg_frame_rate"` // e.g. 30000/1001 = 29.97
}

// frameRate returns the frame rate for the stream.
func (p ProbeOutputStream) frameRate() (float64, error) {
	matches := frameRateMatcher.FindStringSubmatch(p.FrameRate)
	if matches == nil {
		return 0, fmt.Errorf("could not parse ffprobe avg_frame_rate: %s", p.FrameRate)
	}
	a, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		panic(err)
	}
	b, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		panic(err)
	}
	return a / b, nil
}

// ProbeOutput represents the metadata for a video.
type ProbeOutput struct {
	Streams []ProbeOutputStream `json:"streams"`
	Format  struct {
		Duration string `json:"duration"` // e.g. 2769.300000
	} `json:"format"`
}

// firstVideoStream returns the first video stream in the output.
func (p ProbeOutput) firstVideoStream() (ProbeOutputStream, error) {
	for _, stream := range p.Streams {
		if stream.CodecType == "video" {
			return stream, nil
		}
	}
	return ProbeOutputStream{}, fmt.Errorf("no video stream found")
}

// probe retrieves the metadata for the given video file using ffprobe.
func probe(video string) (Metadata, error) {
	raw, _, err := call("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", video)
	if err != nil {
		return Metadata{}, err
	}

	var out ProbeOutput
	err = json.Unmarshal([]byte(raw), &out)
	if err != nil {
		return Metadata{}, err
	}

	secs, err := strconv.ParseFloat(out.Format.Duration, 64)
	if err != nil {
		return Metadata{}, err
	}
	vs, err := out.firstVideoStream()
	if err != nil {
		return Metadata{}, err
	}
	fps, err := vs.frameRate()
	if err != nil {
		return Metadata{}, err
	}

	return Metadata{
		FrameCount:      int64(math.Round(secs * fps)),
		DurationSeconds: secs,
	}, nil
}
