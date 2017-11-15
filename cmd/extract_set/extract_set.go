package main

import (
	"encoding/json"
	"flag"
	"github.com/amarburg/go-frameset"
	"github.com/amarburg/go-lazyquicktime"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {

	// Configuration variables
	var outdir, basedir string
	var doDelete, force bool

	imageDir := "images/"

	flag.StringVar(&outdir, "outdir", ".", "Output directory")
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

	extractor := makeMovieExtractor(set, basedir)

	log.Printf("Extracting %d frames from %s", len(set.Frames), set.Source)

	if set.ImageName == "" {
		log.Print("No image name specified, using default \"image_%06d.png\"")
		set.ImageName = "image_%06d.png"
	}

	pattern := makeImageFilePattern(set.ImageName)

	rootPattern := pattern.SetBaseDir(filepath.Join(outdir, imageDir))
	extractSetFrom(extractor, set.Frames, rootPattern, force, doDelete)

	for name, chunk := range set.Chunks {
		chunkdir := filepath.Join(outdir, name, imageDir)
		chunkPattern := pattern.SetBaseDir(chunkdir)
		extractSetFrom(extractor, chunk.Frames, chunkPattern, force, doDelete)
	}

}

func loadFrameSet(setFile string) FrameSet.FrameSet {
	// Parse the source
	fs, err := os.Open(setFile)
	if err != nil {
		log.Fatalf("Couldn't open set file \"%s\": %s", setFile, err)
	}

	var set FrameSet.FrameSet

	decoder := json.NewDecoder(fs)
	err = decoder.Decode(&set)

	if err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
	}

	return set
}

//
func extractSetFrom(ext lazyquicktime.MovieExtractor, frames []uint64,
	pattern imageFilePattern, force bool, doDelete bool) {

	var existingFiles []string
	if force == false {
		existingFiles = pattern.ExistingFiles()
	}

	for _, frame := range frames {
		outpath := pattern.MakePath(frame)

		var found bool
		before := len(existingFiles)
		existingFiles, found = removeFromSlice(outpath, existingFiles)
		if found == true {
			continue
		}

		img, err := ext.ExtractFrame(frame)

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

	png.Encode(outfile, img)
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
