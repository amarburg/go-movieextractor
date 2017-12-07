package main

import (
	"github.com/amarburg/go-frameset"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type ScrubNail struct {
	Images        []Images
	NumFrames     int
	ScrubnailPath, RelaScrubnailPath string
}

// Some helpful math helpers
func minUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// Some helpful math helpers
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}


func makeScrubNails(im *ImageMaker, chunk frameset.NamedChunk) []ScrubNail {
	// Make configurable later
	framesPerThumb := int(30 * 60)
	framesPerImage := int(30 * 5)

	numThumbs := int(math.Floor(float64(chunk.Range.End-chunk.Range.Start) / float64(framesPerThumb)))

	scrubnails := make([]ScrubNail, numThumbs)

	for i := range scrubnails {

		offset := chunk.Range.Start + uint64(i*framesPerThumb)
		thumbLength := int(minUint64(uint64(framesPerThumb), chunk.Range.End-offset))

		numImages := thumbLength / framesPerImage
		thumbnails := make([]Images, numImages)

		for j := range thumbnails {
			thumbnails[j] = im.MakeImages(offset + uint64(j*framesPerImage))
		}

		// Compose sprites
		spriteSheet := GenerateSpriteSheet(thumbnails)

		scrubnailPath := thumbnails[0].ThumbPath
		scrubnailPath = strings.TrimSuffix(scrubnailPath, filepath.Ext(scrubnailPath))
		scrubnailPath = scrubnailPath + "_scrub.png"

		relaScrubnailPath := thumbnails[0].RelaThumbPath
		relaScrubnailPath = strings.TrimSuffix(relaScrubnailPath, filepath.Ext(relaScrubnailPath))
		relaScrubnailPath = relaScrubnailPath + "_scrub.png"

		out, err := os.Create(scrubnailPath)
		if err != nil {
			log.Printf("Error making scrubnail file %s: %s", scrubnailPath, err)
		} else {
			defer out.Close()
		}

		err = png.Encode(out, spriteSheet)
		if err != nil {
			log.Printf("Error encoding scrubnail file %s: %s", scrubnailPath, err)
		}

		scrubnails[i] = ScrubNail{
			Images:            thumbnails,
			NumFrames:         len(thumbnails),
			RelaScrubnailPath: relaScrubnailPath,
		}

	}

	return scrubnails
}

// Bundle the name of a sprite image with its viewport into the spritesheet.
type Sprite struct {
	image.Rectangle
}

// Compose the spritesheet image and return an array of individual sprite data.
// Based on https://github.com/moovweb/spracker/blob/master/spracker.go
func GenerateSpriteSheet(images []Images) draw.Image {
	var (
		sheetHeight int = 0
		sheetWidth  int = 0
	)

	// Load images
	imgs := make([]image.Image, 0, len(images))
	for _, i := range images {
		imgFile, err := os.Open(i.ThumbPath)
		if err != nil {
			log.Printf("Couldn't open image file '%s'", i.ThumbPath)
			continue
		} else {
			defer imgFile.Close()
		}

		ext := strings.ToLower(filepath.Ext(i.ThumbPath))
		if ext == ".png" {
			img, err := png.Decode(imgFile)
			if err != nil {
				log.Printf("Problem decoding png image in '%s'", i.ThumbPath)
				continue
			}
			imgs = append(imgs, img)
		}
	}

	sprites := make([]Sprite, 0)

	// calculate the size of the spritesheet and accumulate position and padding
	// data for the individual sprites within the sheet
	for _, img := range imgs {
		bounds := img.Bounds()

		newSprite := Sprite{
			image.Rect(sheetWidth, 0, sheetWidth+bounds.Dx(), bounds.Dy()),
		}
		sprites = append(sprites, newSprite)

		sheetWidth += bounds.Dx()
    sheetHeight = maxInt( sheetHeight, bounds.Dy() )

	}

	// create the sheet image
	sheet := image.NewRGBA(image.Rect(0, 0, sheetWidth, sheetHeight))

	// compose the sheet
	for i, img := range imgs {
		draw.Draw(sheet, sprites[i].Rectangle, img, image.Pt(0, 0), draw.Src)
	}

	return sheet
}
