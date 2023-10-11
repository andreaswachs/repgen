/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"os"

	"github.com/andreaswachs/repgen/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "repgen [-i input] [-o output]",
	Short: "repgen is a test report generation tool",
	Long: `repgen generates a single test report HTML file based on output
	from running 'go test -json'.

	repgen reads from stdin by default, but can also read from a specific. Reports are
	written to 'report.html' by default, but can be changed at runtime. See flags for more.
	`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		outputFile, _ := cmd.Flags().GetString("output")
		inputFile, _ := cmd.Flags().GetString("input")

		var inStream io.Reader
		if inputFile == "" {
			inStream = os.Stdin
		} else {
			f, err := os.Open(inputFile)
			if err != nil {
				panic(err)
			}

			defer f.Close()
			inStream = f
		}

		internal.Run(inStream, outputFile)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("input", "i", "", "Input file")
	rootCmd.Flags().StringP("output", "o", "report.html", "Output file")
}
