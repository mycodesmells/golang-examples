package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Verbose bool

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().String("lang", "en", "language to use")

	RootCmd.AddCommand(GreetPlanetCmd)
	RootCmd.AddCommand(GreetTravellerCmd)
}

var RootCmd = &cobra.Command{
	Use: "hello",
	Run: func(cmd *cobra.Command, args []string) {
		if Verbose {
			fmt.Println("About to greet friends from Earth...")
		}
		lang := cmd.Flag("lang").Value.String()
		fmt.Printf("%s world :)\n", greeting(lang))
	},
}
