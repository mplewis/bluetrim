package lib

import (
	"fmt"
	"path"
)

// SuffixFn adds a suffix to a filename but preserves the rest of the path.
func SuffixFn(in string, suffix string) string {
	dir, file := path.Split(in)
	ext := path.Ext(file)
	name := file[:len(file)-len(ext)]
	return path.Join(dir, fmt.Sprintf("%s%s%s", name, suffix, ext))
}

// genOutFn generates an output filename from an input filename.
// Example: /path/to/my/input.mp4 => path/to/my/input_trimmed.mp4
func genOutFn(in string) string {
	return SuffixFn(in, "_trimmed")
}
