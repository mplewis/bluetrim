package lib

import (
	"time"

	"github.com/mplewis/figyr"
)

const desc = "Bluetrim clips dead air from your videos."

// Config defines the runtime app config.
type Config struct {
	Keyframe string        `figyr:"default=start,description=The keyframe to reference as dead air. Possible values are the first frame of the video (\"start\")\\, the last frame (\"end\")\\, or a timestamp in seconds."`
	In       string        `figyr:"required,description=The input file to process."`
	Out      string        `figyr:"optional,description=The destination file for the trimmed video. If unset\\, a filename will be generated from the input file."`
	DryRun   bool          `figyr:"optional,description=Analyze the video but do not generate an output file."`
	Debug    bool          `figyr:"optional,description=Print debug logging information."`
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
