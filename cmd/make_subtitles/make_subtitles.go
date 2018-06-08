package main

import (
	"flag"
	"fmt"
	"github.com/amarburg/go-movieset"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var source, outdir string
	flag.StringVar(&outdir, "outdir", ".", "Directory for resulting subtitle files")

	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("Must specify a multimov .json file on the command line")
	}
	source = flag.Args()[0]

	mm, err := movieset.LoadMultiMov(source)
	if err != nil {
		log.Fatalf("Unable to load MultiMov from \"%s\": %s", source, err)
	}

	// Should look this up in the quicktime file
	dt := 1.0 / 29.97

	offset := uint64(1)

	for _, seq := range mm.Sequence {
		mov, _ := mm.Movies[seq.Hash]
		basename := strings.TrimSuffix(filepath.Base(mov.ShortName), filepath.Ext(mov.ShortName))
		srtFile := filepath.Clean(filepath.Join(outdir, fmt.Sprintf("%s.srt", basename)))

		log.Printf("Writing subtitles to \"%s\"", srtFile)

		file, err := os.Create(srtFile)
		if err != nil {
			log.Fatalf("Unable to open SRT file \"%s\": %s", srtFile, err)
		}

		prevTime := fmtDuration(0.0)

		for i := uint64(0); i < mov.NumFrames; i++ {
			fmt.Fprintf(file, "%d\n", i)
			nextTime := fmtDuration(dt * float64(i+1))
			fmt.Fprintf(file, "%s --> %s\n", prevTime, nextTime)
			fmt.Fprintf(file, "%d\n\n", offset+i)

			prevTime = nextTime
		}

		offset += mov.NumFrames
	}
}

func fmtDuration(d float64) string {

	hours, frac := math.Modf(d / 3600.0)
	mins, frac := math.Modf(frac * 60.0)
	secs, frac := math.Modf(frac * 60.0)

	return fmt.Sprintf("%02d:%02d:%02d,%03d", uint(hours), uint(mins), uint(secs), uint64(frac*1e3))
}
