package frameset

import (
	"fmt"
	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyquicktime"
	"github.com/amarburg/go-multimov"
	"os"
	"path/filepath"
)

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

	ext := filepath.Ext(source)

	switch ext {
	case ".mov":
		fs, err := lazyfs.OpenLocalFile(source)
		if err != nil {
			return nil, fmt.Errorf("Error opening file \"%s\": %s", source, err)
		}

		lqt, err := lazyquicktime.LoadMovMetadata(fs)
		if err != nil {
			return nil, fmt.Errorf("Error parsing Quicktime file \"%s\": %s", source, err)
		}

		return lqt, nil

	case ".json":

		mm, err := multimov.LoadMultiMov(source)
		if err != nil {
			return nil, fmt.Errorf("Error opening multimov file \"%s\": %s", source, err)
		}

		return mm, nil

	default:
		return nil, fmt.Errorf("Unsure what to do with input \"%s\"", source)
	}

}
