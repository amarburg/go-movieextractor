package multimov

import (
	"fmt"
	"github.com/amarburg/go-lazyfs"
	"github.com/amarburg/go-lazyquicktime"
	"path/filepath"
)

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

		mm, err := LoadMultiMov(path)
		if err != nil {
			return nil, err
		}

		return mm, nil

	}

	return nil, fmt.Errorf("Can't make a movie extractor from file %s", path)

}
