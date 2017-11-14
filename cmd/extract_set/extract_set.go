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

	// Configuration variables
	var outdir, basedir string
	var doDelete, force bool

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

	if set.ImageName == "" {
		log.Print("No image name specified, using default \"image_%06d.png\"")
		set.ImageName = "image_%06d.png"
	}

	nameRe := printfToRegexp( set.ImageName )

	var existingFiles []string
	if force == false {
		existingFiles = findImageFiles( outdir, nameRe )
	}

	// Create the source
	ext := filepath.Ext(source)
	set.Source = os.ExpandEnv(set.Source)

	log.Printf("Extracting %d frames from %s", len(set.Frames), set.Source)

	var createdFiles, unmatchedFiles []string

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

		createdFiles,unmatchedFiles = extractSetFrom(lqt, set, outdir, existingFiles)


	case ".json":

		mm, err := multimov.LoadMultiMov(set.Source)
		if err != nil {
			log.Fatalf("Error opening multimov file \"%s\": %s", set.Source, err)
		}

		// If required rewrite the basedir for the source
		if basedir != "" {
			mm.BaseDir = basedir
		}

	createdFiles,unmatchedFiles = extractSetFrom(mm, set, outdir, existingFiles)

	default:
		log.Fatalf("Unsure what to do with input \"%s\"", source)
	}



	log.Printf("Created %d image files, %d orphaned files", len(createdFiles), len(unmatchedFiles))

	if doDelete {
		log.Printf("Deleting orphaned image files")

		for _, filename := range unmatchedFiles {
			fullpath := filepath.Clean(filepath.Join(outdir,filename))
			log.Printf("Deleting file %s", fullpath)
			err := os.Remove(fullpath)
			if err != nil {
				log.Printf("Couldn't delete \"%s\": %s", fullpath, err)
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
func extractSetFrom(ext lazyquicktime.MovieExtractor, set frameSet,
										outdir string, existing []string) (extractedFiles, unmatchedFiles []string) {

	extractedFiles = make([]string, 0, len(set.Frames))

	for _, frame := range set.Frames {
		outname := fmt.Sprintf(set.ImageName, frame)

		var found bool
		existing,found = removeFromSlice(outname, existing)
		if found == true {
			log.Printf("File \"%s\" exists, skipping", outname)
			continue
		}

		img, err := ext.ExtractFrame(frame)

		if err != nil {
			log.Fatalf("Unable to extract frame: %s", err)
		}

		outpath := filepath.Clean(filepath.Join(outdir,outname))
		writeImage(img, outpath)

		extractedFiles = append(extractedFiles, outname)
	}

	return extractedFiles, existing
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

func removeFromSlice(a string, list []string) (out []string, found bool) {
	for i, b := range list {
		if b == a {
			list[len(list)-1], list[i] = list[i], list[len(list)-1]
			return list[:len(list)-1], true
		}
	}
	return list, false
}


func printfToRegexp( printfFmt string ) (* regexp.Regexp) {
	// Find existing files

	digitsRe,_ := regexp.Compile("%[0-9]*d")

	nameRegexp := digitsRe.ReplaceAllString(printfFmt, "[\\d]*")
	log.Printf("Converted filename pattern \"%s\" to regex \"%s\"", printfFmt, nameRegexp)

	nameRe, err := regexp.Compile(nameRegexp)

	if( err != nil ){
		log.Fatalf("Unable to compile filename regexp \"%s\"", nameRegexp)
	}

	return nameRe
}



func findImageFiles( path string, nameRe *regexp.Regexp ) ([]string) {
	dir, _ := os.Open(filepath.Dir(path))
	defer dir.Close()

	files, _ := dir.Readdirnames(0)
	existing := make([]string,0,len(files))

	for _, filename := range files {
		//log.Printf("Checking %s", filename)

		if nameRe.MatchString(filename) {
			//log.Printf("File %s matches pattern", filename)
			existing = append(existing,filename)
		}
	}

	return existing
}
