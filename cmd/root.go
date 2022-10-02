/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/jdfergason/jflix-split/scanner"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	show    string
	season  int
	episode int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jflix-split <video file>",
	Args:  cobra.ExactArgs(1),
	Short: "split a recording into multiple episodes",
	Long: `jflix-split attempts to identify the beginning and end of individual episodes
	or movies in a recording and split them into individual files.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Str("File", args[0]).Msg("Splitting video file into constituents")
		scan := scanner.NewScanner(args[0], show, season, episode)
		scan.FindSegments()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&show, "show", "TV Show", "show name")
	rootCmd.Flags().IntVarP(&season, "season", "s", 1, "season of first item in video input")
	rootCmd.Flags().IntVarP(&episode, "episode", "e", 1, "starting episode of first item in video input")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
