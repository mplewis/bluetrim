package lib

import (
	"fmt"
	"path"
	"time"

	"github.com/mplewis/figyr"
)

const desc = "Bluetrim clips dead air from your videos."

// Config defines the runtime app config.
type Config struct {
	Keyframe string        `figyr:"default=start,description=The keyframe to use for dead air. Possible values are the first frame of the video (start)\\, the last frame (end)\\, a frame number\\, or any valid timecode (00:00:00.000)."`
	In       string        `figyr:"required,description=The input file to process."`
	Out      string        `figyr:"optional,description=The destination file for the trimmed video. If unset\\, a filename will be generated from the input file."`
	DryRun   bool          `figyr:"optional,description=If set\\, the program will analyze the video but will not generate an output file."`
	Interval time.Duration `figyr:"default=1m,description=The interval between frames to use to check for dead air."`
}

// LoadConfig loads the runtime app config.
func LoadConfig() Config {
	var cfg Config
	figyr.New(desc).MustParse(&cfg)
	if cfg.DryRun && cfg.Out != "" {
		panic("arguments Out and DryRun cannot be used simultaneously")
	}
	if cfg.Out == "" {
		cfg.Out = genOutFn(cfg.In)
	}
	return cfg
}

// genOutFn generates an output filename from an input filename.
// Example: /path/to/my/input.mp4 => path/to/my/input_trimmed.mp4
func genOutFn(in string) string {
	dir, file := path.Split(in)
	ext := path.Ext(file)
	name := file[:len(file)-len(ext)]
	return path.Join(dir, fmt.Sprintf("%s_trimmed%s", name, ext))
}
