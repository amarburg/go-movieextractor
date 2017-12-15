package main

import (
	"net/http"
	"github.com/spf13/cobra"
	"log"
)

var cmdHttpServer = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	Run:   runHttpServer,
}

func init() {
	//cmdProofSheet.PersistentFlags().StringVar(&templateFile, "template", filepath.Join( DefaultTemplateDir(), "synopsis.tmpl"), "Template file to use")
	//cmdProofSheet.PersistentFlags().Uint64Var(&step, "step", uint64(math.Trunc(29.97*60)), "")
}

func runHttpServer(cmd *cobra.Command, args []string) {

	outdir := OutputDir

	log.Printf("Listening on localhost:%04d", 9080 )

	http.Handle("/", http.FileServer(http.Dir(outdir)))
	log.Fatal(http.ListenAndServe(":9080", nil))

	// source := args[0]
  //
	// set, err := frameset.LoadFrameSet(source)
	// if err != nil {
	// 	log.Fatalf("Unable to load FrameSet from \"%s\": %s", source, err)
	// }
  //
	// // Construct multimov from image set
	// multimovPath := os.ExpandEnv(set.Source)
	// mm, err := multimov.LoadMultiMov(multimovPath)
	// if err != nil {
	// 	log.Fatalf("Unable to load MultiMov from \"%s\": %s", multimovPath, err)
	// }
  //
	// set.NumFrames = mm.NumFrames()
  //
	// outTree := NewDirOutTree(outdir)
	// im := NewImageMaker(mm, outTree)
  //
	// indexfile := outTree.join("index.html")
	// outfile, err := os.Create(indexfile)
	// if err != nil {
	// 	log.Fatalf("Unable to open the output file \"%s\": %s", indexfile, outfile)
	// }
	// defer outfile.Close()
  //
	// // Funtion map
	// fmap := template.FuncMap{
	// 	"makeImages": im.MakeImages,
	// 	"makeScrubNails": func(chunk frameset.Chunk) []ScrubNail {
	// 		return makeScrubNails(im, chunk)
	// 	},
	// 	"frameSetName": func(set frameset.FrameSet) string {
	// 		return filepath.Base(set.Source)
	// 	},
	// }
  //
	// t := template.New("synopsis.tmpl")
	// t = t.Funcs(fmap)
  //
	// // Load layouts
	// t = template.Must(t.ParseGlob(filepath.Join(DefaultTemplateDir(), "layouts/*.tmpl")))
	// t = template.Must(t.ParseFiles(filepath.Join(DefaultTemplateDir(), "synopsis.tmpl")))
  //
	// err = t.Execute(outfile, set)
	// if err != nil {
	// 	panic(err)
	// }

}
