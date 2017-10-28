package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var GreetTravellerCmd = &cobra.Command{
	Use: "traveller",
	Run: func(cmd *cobra.Command, args []string) {
		lang := cmd.Flag("lang").Value.String()
		greet := greeting(lang)

		if len(args) == 0 {
			fmt.Printf("%s travellers!\n", greet)
			os.Exit(0)
		}

		fmt.Printf("%s travellers from %s!\n", greet, args[0])
	},
}
