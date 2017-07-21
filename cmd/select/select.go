package main


import (
  "flag"
  "log"
  "encoding/json"
  "os"
  "fmt"
  "image"
  "image/png"
  "path/filepath"
)


type FrameSource interface {
  Get(uint) (image.Image, error)
}



type ImageRegexpSource struct {
  re string
}


func NewImageRegexpSource( re string ) (ImageRegexpSource) {
  return ImageRegexpSource{
    re: re,
  }
}

func (src ImageRegexpSource) Get( frame uint) (image.Image, error ) {

  fileGlob := fmt.Sprintf(src.re, frame)

  log.Printf("Looking for %s", fileGlob)

  matches,err := filepath.Glob( fileGlob )

  if err == filepath.ErrBadPattern {
    return nil, fmt.Errorf("Badly formatted file regexp: %s", fileGlob )
  } else if len(matches) == 0 {
    log.Printf("Found no files which match \"%s\"", fileGlob )
    return nil, err
  } else if len(matches) != 1 {
    log.Printf("File glob \"%s\" didn't resolve to just one file, I found: %s", fileGlob, matches)
    return nil, err
  }

  fileName := matches[0]

  f, err := os.Open( fileName )

  if err != nil  {
    return nil, err
}

  img, fmt, err := image.Decode( f )

  if err != nil {
    log.Printf("Loaded image %s of type %s", fileName, fmt)
}

  return img, err

}




type Config struct {
  jsonFile, outRe string

  force, delete bool

  source  FrameSource
}

type JsonData struct {
  Frames []uint
}


func ParseFlags() (Config,error) {



  var forceArg = flag.Bool("force", false, "Force creation of images")
  var deleteArg = flag.Bool("delete", false, "Delete existing file which aren't in the JSON file")
  var imageRe = flag.String("image-re","","Image regexp")
  var outputRe = flag.String("output-re", "image_%06d.png", "Output directory")

  flag.Parse()

  var conf = Config{
    force:    *forceArg,
    delete:   *deleteArg,
    jsonFile: flag.Arg(0),
  }

  if conf.jsonFile == "" {
    return conf, fmt.Errorf("Must supply file on command line")
  }

  conf.outRe = *outputRe


  fmt.Println(*imageRe)
  if *imageRe != "" {
    conf.source = NewImageRegexpSource( *imageRe )
  } else {
    log.Fatal("No image source specified")
}



  return conf, nil
}



func main() {

  conf,err := ParseFlags()

  if err != nil {
      log.Fatal("Error parsing args: ", err)
  }

  f,err := os.Open( conf.jsonFile )
  defer f.Close()

  if err != nil {
    log.Fatal("Error opening file \"", conf.jsonFile, "\":", err )
  }

  decoder := json.NewDecoder( f )

  var jsonFile = JsonData{}

  err = decoder.Decode( &jsonFile )
  if err != nil {
    log.Fatalf("Error parsing json: %s", err )
  }

  fmt.Printf("Parsed file %s with %d frames\n", conf.jsonFile, len(jsonFile.Frames) )

  for _,i := range jsonFile.Frames {


    outpath := fmt.Sprintf( conf.outRe, i )
    if _, err = os.Stat(outpath); !os.IsNotExist(err) && !conf.force {
      log.Printf("File \"%s\" exists, skipping", outpath)
      continue
    }


    img,err := conf.source.Get(i)

    if err != nil {
      log.Printf("Couldn't load frame number %d: %s", i, err )
      continue
    }

    out,err := os.Create(outpath)
    if err != nil {
  log.Fatalf("Could not open output file %s", err)

}
  defer out.Close()

  png.Encode( out, img )

  log.Printf("Saved to %s", outpath)
}

}
