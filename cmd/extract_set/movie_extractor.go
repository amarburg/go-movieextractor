package main

import (
	"github.com/amarburg/go-frameset"
	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyquicktime"
	"github.com/amarburg/go-multimov"
	"log"
	"os"
	"path/filepath"
)

func makeMovieExtractor(set FrameSet.FrameSet, basedir string) lazyquicktime.MovieExtractor {
	// Create the source
	source := os.ExpandEnv(set.Source)
	ext := filepath.Ext(source)

	log.Printf("Extracting %d frames from %s", len(set.Frames), source)

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

		return lqt

	case ".json":

		mm, err := multimov.LoadMultiMov(source)
		if err != nil {
			log.Fatalf("Error opening multimov file \"%s\": %s", source, err)
		}

		// If required rewrite the basedir for the source
		if basedir != "" {
			mm.BaseDir = basedir
		}

		return mm

	default:
		log.Fatalf("Unsure what to do with input \"%s\"", source)
	}
	return nil
}
