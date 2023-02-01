package main

import (
	"fmt"

	"github.com/mplewis/bluetrim/lib"
)

// main runs the program.
func main() {
	cfg := lib.LoadConfig()
	fmt.Println(lib.Probe(cfg.In))
	fmt.Println(lib.CmpImages("/Users/mplewis/tmp/framex/output_0001.jpg", "/Users/mplewis/tmp/framex/output_0002.jpg"))
}
