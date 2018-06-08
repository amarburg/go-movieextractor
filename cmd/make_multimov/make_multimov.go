package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyquicktime"
	"github.com/amarburg/go-movieset"
	"os"
	"path/filepath"
)

func main() {
	outfileFlag := flag.String("output", "multimov.json", "Name of output JSON file")
	basedirFlag := flag.String("basedir", "", "Name of basedir to use")

	flag.Parse()

	outfile := *outfileFlag

	var basedir string
	if *basedirFlag != "" {
		basedir = *basedirFlag
	} else {
		basedir = filepath.Dir(outfile)
	}

	fmt.Printf("Outfile:  %s\n", outfile)
	fmt.Printf("Basedir:  %s\n", basedir)

	// Create the Multimov
	mm := movieset.NewMultiMov()

	for _, pathStr := range flag.Args() {

		matches, _ := filepath.Glob(pathStr)

		for _, path := range matches {

			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Printf("Path \"%s\" doesn't exist\n", path)
				continue
			}

			lqt, err := makeLazyQuicktime(path)

			if err != nil {
				fmt.Println("Error making lazyquicktime:", err)
				continue
			}

			mr := movieset.MovRecordFromLqt(lqt)

			// Rebase to the specified basedir
			mr.Relapath, _ = filepath.Rel(basedir, path)

			mm.AppendMovie(mr)
		}
	}

	file, _ := os.Create(outfile)
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(mm)
}

func makeLazyQuicktime(path string) (*lazyquicktime.LazyQuicktime, error) {

	fs, err := lazyfs.OpenLocalFile(path)
	if err != nil {
		return nil, err
	}

	lqt, err := lazyquicktime.LoadMovMetadata(fs)
	if err != nil {
		return nil, err
	}

	return lqt, err
}
