package main

import (
	// "fmt"
	// "os"

	//homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

var OutputDir string

// func init() {
//
//   // rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
//   // rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
//   // rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
//   // rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
//   // rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
//   // viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
//   // viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
//   // viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
//   // viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
//   // viper.SetDefault("license", "apache")
// }

func initConfig() {
	// // Don't forget to read config either from cfgFile or from home directory!
	// if cfgFile != "" {
	//   // Use config file from the flag.
	//   viper.SetConfigFile(cfgFile)
	// } else {
	//   // Find home directory.
	//   home, err := homedir.Dir()
	//   if err != nil {
	//     fmt.Println(err)
	//     os.Exit(1)
	//   }
	//
	//   // Search config in home directory with name ".cobra" (without extension).
	//   viper.AddConfigPath(home)
	//   viper.SetConfigName(".cobra")
	// }

	// if err := viper.ReadInConfig(); err != nil {
	//   fmt.Println("Can't read config:", err)
	//   os.Exit(1)
	// }
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "make_proof_sheet",
		Short: "Hugo is a very fast static site generator",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}

	//cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&OutputDir, "output", "./_html", "Output directory (default is ./_html)")

	rootCmd.AddCommand(cmdProofSheet)
	rootCmd.AddCommand(cmdSynopsis)
	rootCmd.Execute()
}

// package main
//
// import (
// 	"flag"
// 	"github.com/google/subcommands"
// 	"context"
// )
//
// var count uint
//
// func main() {
//
// 	topLevelFlags := flag.NewFlagSet("top level", flag.ExitOnError)
//
// 	var outdir string
// 	flag.StringVar(&outdir, "outdir", "_html", "Directory for resulting html files")
//
// 	commander := subcommands.NewCommander( topLevelFlags, "make_proof_sheet" )
//
// 	commander.Register( subcommands.HelpCommand(), "" )
// 	commander.Register( subcommands.FlagsCommand(), "" )
// 	commander.Register( NewProofSheetCommand(), "" )
// 	commander.Register( NewSynopsisCommand(), "" )
//
// 	commander.Execute( context.Background(), outdir )
//
// }
