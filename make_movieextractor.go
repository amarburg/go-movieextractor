package frameset

import (
	"fmt"
	"github.com/amarburg/go-lazyquicktime"
	"github.com/amarburg/go-multimov"
	"github.com/amarburg/go-lazyfs"
	"os"
	"path/filepath"
)


// Technically this could be in multimov?
func MovieExtractorFromPath(path string) (lazyquicktime.MovieExtractor, error) {

	if filepath.Ext(path) == ".mov" {

		file, err := lazyfs.SourceFromPath(path)

		if err != nil {
			return nil, err
		}

		qtInfo, err := lazyquicktime.LoadMovMetadata(file)

		if err != nil {
			return nil, err
		}

		return qtInfo, nil

	} else if filepath.Ext(path) == ".json" {

		fmt.Println(path)

		mm, err := multimov.LoadMultiMov(path)
		if err != nil {
			return nil, err
		}

		return mm, nil

	}

	return nil, fmt.Errorf("Can't make a movie extractor from file %s", path)

}


func (set FrameSet) MovieExtractor() (lazyquicktime.MovieExtractor, error) {
	// Create the source
	source := os.ExpandEnv(set.Source)

	// If path is relative, evaluate it relative to the input file...
	if !filepath.IsAbs(set.Source) {
		if len(set.filepath) == 0 {
			return nil, fmt.Errorf("Source path in frameset is relative, but don't know filepath for frameset")
		}

		source = filepath.Clean(filepath.Join(filepath.Dir(set.filepath), source))
	}

	return MovieExtractorFromPath( source )


}
