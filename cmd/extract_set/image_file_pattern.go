package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type imageFilePattern struct {
	printf string
	re     *regexp.Regexp
	dir    string
}

func makeImageFilePattern(printfPattern string) imageFilePattern {
	return imageFilePattern{
		printf: printfPattern,
		re:     printfToRegexp(printfPattern),
	}
}

func printfToRegexp(printfFmt string) *regexp.Regexp {
	// Find existing files

	digitsRe, _ := regexp.Compile("%[0-9]*d")

	nameRegexp := digitsRe.ReplaceAllString(printfFmt, "[\\d]*")
	log.Printf("Converted filename pattern \"%s\" to regex \"%s\"", printfFmt, nameRegexp)

	nameRe, err := regexp.Compile(nameRegexp)

	if err != nil {
		log.Fatalf("Unable to compile filename regexp \"%s\"", nameRegexp)
	}

	return nameRe
}

func (ifp imageFilePattern) SetBaseDir(dir string) imageFilePattern {
	return imageFilePattern{
		printf: ifp.printf,
		re:     ifp.re,
		dir:    dir,
	}
}

func (ifp imageFilePattern) ExistingFiles() []string {
	fmt.Printf("Checking directory \"%s\"", ifp.dir)
	dir, _ := os.Open(ifp.dir)
	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		log.Printf("Error reading filenames: %s", err)
	}
	existing := make([]string, 0, len(files))

	for _, filename := range files {
		log.Printf("Checking %s", filename)

		if ifp.re.MatchString(filename) {
			log.Printf("File %s matches pattern", filename)
			existing = append(existing, filepath.Join(ifp.dir, filename))
		} else {
			log.Printf("File %s does not match pattern", filename)
		}
	}

	return existing
}

func (ifp imageFilePattern) MakePath(frame uint64) string {
	return filepath.Join(ifp.dir, fmt.Sprintf(ifp.printf, frame))
}
