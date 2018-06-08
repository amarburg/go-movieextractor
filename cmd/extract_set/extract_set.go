package main

import (
	"encoding/json"
	"flag"
	"github.com/amarburg/go-movieset"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {

	// Configuration variables
	var outdir, basedir string
	var doDelete, force bool

	imageDir := "images/"

	flag.StringVar(&outdir, "outdir", "", "Output directory")
	flag.StringVar(&basedir, "basedir", "", "Base directory")
	flag.BoolVar(&doDelete, "delete", false, "Delete existing files not specified in set file")
	flag.BoolVar(&force, "force", false, "Replace existing image files")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("Need to specify a setfile on command line")
	}

	source := flag.Args()[0]

	// Load FrameSet JSON file
	set := loadFrameSet(source)
	extractor, err := set.MovieExtractor()

	if err != nil {
		log.Fatalf("Unable to open source \"%s\"", source)
	}

	log.Printf("Extracting %d frames from %s", set.NumFrames, set.Source)

	if set.ImageName == "" {
		set.ImageName = "image_%06d.png"
		log.Printf("No image name specified, using default \"%s\"\n", set.ImageName)
	}

	if outdir == "" {
		outdir = filepath.Dir(source)
	}

	pattern := makeImageFilePattern(set.ImageName)

	_ = os.MkdirAll(filepath.Join(outdir, imageDir), os.ModePerm)
	// rootPattern := pattern.SetBaseDir(filepath.Join(outdir, imageDir))
	// extractSetFrom(extractor, set.Frames, rootPattern, force, doDelete)

	for _, chunk := range set.Chunks {
		chunkdir := filepath.Join(outdir, chunk.Name, imageDir)
		_ = os.MkdirAll(chunkdir, os.ModePerm)

		chunkPattern := pattern.SetBaseDir(chunkdir)
		extractSetFrom(extractor, chunk.Frames, chunkPattern, force, doDelete)
	}

}

func loadFrameSet(setFile string) movieset.FrameSet {
	// Parse the source
	fs, err := os.Open(setFile)
	if err != nil {
		log.Fatalf("Couldn't open set file \"%s\": %s", setFile, err)
	}

	var set movieset.FrameSet

	decoder := json.NewDecoder(fs)
	err = decoder.Decode(&set)

	if err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
	}

	return set
}

//
func extractSetFrom(ext movieset.MovieExtractor, frames []uint64,
	pattern imageFilePattern, force bool, doDelete bool) {

	var existingFiles []string
	if !force {
		existingFiles = pattern.ExistingFiles()
	}

	for _, frame := range frames {
		outpath := pattern.MakePath(frame)

		var found bool
		existingFiles, found = removeFromSlice(outpath, existingFiles)
		if found {
			continue
		}

		start := time.Now()
		img, err := ext.ExtractFrame(frame)
		log.Printf("Extraction took %s", time.Since(start).String())

		if err != nil {
			log.Fatalf("Unable to extract frame: %s", err)
		}

		log.Printf("Saving image to \"%s\"", outpath)
		writeImage(img, outpath)
	}

	if doDelete {
		log.Printf("Deleting %d orphaned image files", len(existingFiles))

		for _, filename := range existingFiles {
			log.Printf("Deleting file %s", filename)
			err := os.Remove(filename)
			if err != nil {
				log.Printf("Couldn't delete \"%s\": %s", filename, err)
			}
		}
	}
}

func writeImage(img image.Image, path string) {
	outfile, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error creating image file \"%s\": %s", path, err)
	}

	_ =png.Encode(outfile, img)
}

func removeFromSlice(a string, list []string) (out []string, found bool) {
	for i, b := range list {
		if b == a {
			list[len(list)-1], list[i] = list[i], list[len(list)-1]
			return list[:len(list)-1], true
		}
	}
	return list, false
}
