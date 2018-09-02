package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MichaelTJones/walk"
	"github.com/ameliaikeda/thumbnailler/processor"
	bimg "gopkg.in/h2non/bimg.v1"
)

var (
	source, destination    string
	width, height, quality int

	del    bool
	search string
)

func init() {
	flag.StringVar(&source, "src", ".", "the source directory where images reside.")
	flag.StringVar(&destination, "dst", "", "the destination for thumbnails")
	flag.IntVar(&width, "width", 300, "the width of generated thumbnails. aspect ratio is preserved.")
	flag.IntVar(&height, "height", 300, "the height of generated thumbnails. aspect ratio is preserved.")
	flag.IntVar(&quality, "quality", 80, "the defined JPEG quality for output.")
	flag.BoolVar(&del, "delete", false, "switches to delete mode on a search")
	flag.StringVar(&search, "search", "", "a path match (non-regex) to use for fast file deletion")

	flag.Parse()
}

func main() {
	if del {
		if search == "" {
			fmt.Fprintln(os.Stderr, "error: search is required when deleting\n")
			flag.Usage()

			os.Exit(2)
		}

		processor.SetSearch(search)

		err := walk.Walk(processor.Expand(source), processor.Delete)
		if err != nil {
			fmt.Fprintf(os.Stderr, "walk error: %+v", err)

			os.Exit(1)
		}

		return
	}

	if destination == "" {
		fmt.Fprintln(os.Stderr, "error: dst must be set\n")
		flag.Usage()

		os.Exit(2)
	}

	processor.SetDestination(destination)
	processor.SetWidth(width)
	processor.SetHeight(height)
	processor.SetQuality(quality)

	defer bimg.Shutdown()

	err := walk.Walk(processor.Expand(source), processor.Process)
	if err != nil {
		fmt.Fprintf(os.Stderr, "walk error: %+v", err)

		os.Exit(1)
	}
}
