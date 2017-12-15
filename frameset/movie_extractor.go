package frameset

import (
	"fmt"
	"github.com/amarburg/go-frameset/multimov"
	"github.com/amarburg/go-lazyquicktime"
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

	return multimov.MovieExtractorFromPath(source)

}
