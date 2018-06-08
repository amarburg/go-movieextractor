package main

import (
	"flag"
	//"fmt"
	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyquicktime"
	// "github.com/amarburg/go-quicktime"
	"github.com/amarburg/go-movieset"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {

	var frame uint64
	var outfile, basedir string
	flag.Uint64Var(&frame, "frame", 1, "Frame number to extract")
	flag.StringVar(&outfile, "outfile", "image.png", "outfile image name")
	flag.StringVar(&basedir, "basedir", "", "Base directory for multimov file")

	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("Need to specify either a movie file or multimov file on the command line")
	}

	source := flag.Args()[0]
	ext := filepath.Ext(source)

	log.Printf("Extracting frame %d from %s", frame, source)

	switch ext {
	case ".mov":
		fs, err := lazyfs.OpenLocalFile(source)
		if err != nil {
			log.Fatalf("Error opening file \"%s\": %s", source, err)
		}

		lqt, err := lazyquicktime.LoadMovMetadata(fs)
		if err != nil {
			log.Fatalf("Error parsing Quicktime file \"%s\": %s", source, err)
		}

		extractAndSave(lqt, frame, outfile)

	case ".json":

		mm, err := movieset.LoadMultiMov(source)
		if err != nil {
			log.Fatalf("Error opening multimov file \"%s\": %s", source, err)
		}

		if basedir != "" {
			mm.BaseDir = basedir
		}

		extractAndSave(mm, frame, outfile)

	default:
		log.Fatalf("Unsure what to do with input \"%s\"", source)
	}

}

func extractAndSave(ext movieset.MovieExtractor, frame uint64, outfile string) {
	img, err := ext.ExtractFrame(frame)

	if err != nil {
		log.Fatalf("Unable to extract frame: %s", err)
	}

	writeImage(img, outfile)
}

func writeImage(img image.Image, path string) {
	outfile, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error opening image file \"%s\": %s", path, err)
	}

	err = png.Encode(outfile, img)
	if err != nil {
		log.Fatalf("Error encoding image file \"%s\": %s", path, err)
	}
}
