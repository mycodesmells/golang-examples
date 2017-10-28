package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: `Summary of cache contents`,
	Long:  `Displays a short summary of what is currently cached`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			cacheKey := fmt.Sprintf("cache:clients:%s", args[0])
			ks, err := redisClient.Keys(cacheKey).Result()
			if err != nil {
				fmt.Println("Failed to list clients cache")
				os.Exit(1)
			}
			if len(ks) == 0 {
				fmt.Printf("Client %s is not present in cache\n", args[0])
				os.Exit(0)
			}
			fmt.Printf("Client %s is present in cache\n", args[0])
			os.Exit(0)
		}

		ks, _ := redisClient.Keys("cache:clients:*").Result()
		fmt.Println("Current cache summary:")
		fmt.Printf("‣ clients cached: %d\n", len(ks))

		if Verbose {
			for _, k := range ks {
				siteID := strings.Replace(k, "cache:clients:", "", -1)
				fmt.Printf("\t‣ %s\n", siteID)
			}
		}

	},
}
