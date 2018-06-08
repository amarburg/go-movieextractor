package main

import (
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var cmdHttpServer = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	Run:   runHttpServer,
}

//func init() {}

func runHttpServer(cmd *cobra.Command, args []string) {

	outdir := OutputDir

	log.Printf("Listening on localhost:%04d", 9080)

	http.Handle("/", http.FileServer(http.Dir(outdir)))
	log.Fatal(http.ListenAndServe(":9080", nil))
}
