package main

import (
	//"flag"
	"github.com/amarburg/go-movieset"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

var cmdSynopsis = &cobra.Command{
	Use:   "synopsis [image set file]",
	Short: "Generates a synopsis page sheet from an image set",
	Long:  "a long string explaining the command and giving usage information",
	Args:  cobra.MinimumNArgs(1),
	Run:   runSynopsis,
}

func init() {
	//cmdProofSheet.PersistentFlags().StringVar(&templateFile, "template", filepath.Join( DefaultTemplateDir(), "synopsis.tmpl"), "Template file to use")
	//cmdProofSheet.PersistentFlags().Uint64Var(&step, "step", uint64(math.Trunc(29.97*60)), "")
}

func runSynopsis(cmd *cobra.Command, args []string) {

	outdir := OutputDir

	source := args[0]

	set, err := movieset.LoadFrameSet(source)
	if err != nil {
		log.Fatalf("Unable to load FrameSet from \"%s\": %s", source, err)
	}

	// Construct multimov from image set
	multimovPath := os.ExpandEnv(set.Source)
	mm, err := movieset.LoadMultiMov(multimovPath)
	if err != nil {
		log.Fatalf("Unable to load MultiMov from \"%s\": %s", multimovPath, err)
	}

	set.NumFrames = mm.NumFrames()

	outTree := NewDirOutTree(outdir)
	im := NewImageMaker(mm, outTree)

	indexfile := outTree.join("index.html")
	outfile, err := os.Create(indexfile)
	if err != nil {
		log.Fatalf("Unable to open the output file \"%s\": %s", indexfile, err.Error())
	}
	defer outfile.Close()

	// Funtion map
	fmap := template.FuncMap{
		"makeImages": im.MakeImages,
		"makeScrubNails": func(chunk movieset.Chunk) []ScrubNail {
			return makeScrubNails(im, chunk)
		},
		"frameSetName": func(set movieset.FrameSet) string {
			return filepath.Base(set.Source)
		},
	}

	t := template.New("synopsis.tmpl")
	t = t.Funcs(fmap)

	// Load layouts
	t = template.Must(t.ParseGlob(filepath.Join(DefaultTemplateDir(), "layouts/*.tmpl")))
	t = template.Must(t.ParseFiles(filepath.Join(DefaultTemplateDir(), "synopsis.tmpl")))

	err = t.Execute(outfile, set)
	if err != nil {
		panic(err)
	}

}
