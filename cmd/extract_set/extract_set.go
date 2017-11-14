package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyquicktime"
	"github.com/amarburg/go-multimov"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type frameSet struct {
	Source    string
	Frames    []uint64
	ImageName string
}

func main() {

	var outdir, basedir string
	var doDelete bool
	flag.StringVar(&outdir, "outdir", ".", "Output directory")
	flag.StringVar(&basedir, "basedir", "", "Base directory")

	flag.BoolVar(&doDelete, "delete", false, "Delete existing files not specified in set file")

	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("Need to specify a setfile on command line")
	}

	source := flag.Args()[0]

	set := loadFrameSet(source)

	if set.ImageName == "" {
		log.Print("No image name specified, using default \"image_%06d.png\"")
		set.ImageName = "image_%06d.png"
	}

	// Find existing files
	transRe, _ := regexp.Compile("%[0-9]*d")

	namePattern := transRe.ReplaceAllString(set.ImageName, "[\\d]*")
	log.Printf("Converted filename pattern \"%s\" to regex \"%s\"", set.ImageName, namePattern)

	nameRe, _ := regexp.Compile(namePattern)
	dir, _ := os.Open(filepath.Dir(outdir))
	defer dir.Close()

	files, _ := dir.Readdirnames(0)
	existingFiles = make([]string,0,len(files))

	for _, filename := range files {
		//log.Printf("Checking %s", filename)

		if nameRe.MatchString(filename) {
			//log.Printf("File %s matches pattern", filename)
			filename = append(filename,existingFiles)
		}
	}



	// Create the source
	ext := filepath.Ext(source)
	set.Source = os.ExpandEnv(set.Source)

	log.Printf("Extracting %d frames from %s", len(set.Frames), set.Source)

	var createdFiles []string
	switch ext {
	case ".mov":
		fs, err := lazyfs.OpenLocalFile(set.Source)
		if err != nil {
			log.Fatalf("Error opening file \"%s\": %s", set.Source, err)
		}

		lqt, err := lazyquicktime.LoadMovMetadata(fs)
		if err != nil {
			log.Fatalf("Error parsing Quicktime file \"%s\": %s", set.Source, err)
		}

		createdFiles = extractSetFrom(lqt, set, outdir)

	case ".json":

		mm, err := multimov.LoadMultiMov(set.Source)
		if err != nil {
			log.Fatalf("Error opening multimov file \"%s\": %s", set.Source, err)
		}

		// If required rewrite the basedir for the source
		if basedir != "" {
			mm.BaseDir = basedir
		}

		createdFiles = extractSetFrom(mm, set, outdir)

	default:
		log.Fatalf("Unsure what to do with input \"%s\"", source)
	}

	log.Printf("Created %d image files", len(createdFiles))

	if doDelete {
		log.Printf("Handling doDelete")

		// Convert set.ImageName to a regex
		transRe, _ := regexp.Compile("%[0-9]*d")

		namePattern := transRe.ReplaceAllString(set.ImageName, "[\\d]*")
		log.Printf("Converted filename pattern \"%s\" to regex \"%s\"", set.ImageName, namePattern)

		nameRe, _ := regexp.Compile(namePattern)

		dir, _ := os.Open(filepath.Dir(outdir))

		files, _ := dir.Readdirnames(0)
		for _, filename := range files {
			//log.Printf("Checking %s", filename)

			if nameRe.MatchString(filename) {
				//log.Printf("File %s matches pattern", filename)

				if !stringInSlice(filename, createdFiles) {
					log.Printf("Deleting file %s", filename)
					os.Remove(filename)
				}
			}
		}
	}

}

func loadFrameSet(setFile string) frameSet {
	// Parse the source
	fs, err := os.Open(setFile)
	if err != nil {
		log.Fatalf("Couldn't open set file \"%s\": %s", setFile, err)
	}

	var set frameSet

	decoder := json.NewDecoder(fs)
	err = decoder.Decode(&set)

	if err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
	}

	return set
}

//
func extractSetFrom(ext lazyquicktime.MovieExtractor, set frameSet, outdir string) []string {

	extractedFiles := make([]string, 0, len(set.Frames))

	for _, frame := range set.Frames {
		img, err := ext.ExtractFrame(frame)

		if err != nil {
			log.Fatalf("Unable to extract frame: %s", err)
		}

		outfile := filepath.Clean(filepath.Join(outdir, fmt.Sprintf(set.ImageName, frame)))

		writeImage(img, outfile)

		extractedFiles = append(extractedFiles, outfile)
	}

	return extractedFiles
}

func writeImage(img image.Image, path string) {
	outfile, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error creating image file \"%s\": %s", path, err)
	}

	png.Encode(outfile, img)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
