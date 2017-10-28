package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var client string

func init() {
	clearCmd.Flags().StringVarP(&client, "client", "c", "", "client ID")
	clearCmd.AddCommand(clearAllCmd)
	clearCmd.MarkFlagRequired("client")
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: `Clears cache contents`,
	Long:  `Clears cached data about clients`,
	Run: func(cmd *cobra.Command, args []string) {
		clientID := cmd.Flag("client").Value.String()
		cacheKey := fmt.Sprintf("cache:clients:%s", clientID)

		fmt.Printf("Clearing cached item of %s... ", clientID)
		err := redisClient.Del(cacheKey).Err()
		if err != nil {
			fmt.Println("✘")
			os.Exit(1)
		}
		fmt.Println("✔")
	},
}

var clearAllCmd = &cobra.Command{
	Use:   "all",
	Short: `Clears cache contents`,
	Long:  `Clears cached data about clients`,
	Run: func(cmd *cobra.Command, args []string) {
		ks, _ := redisClient.Keys("cache:clients:*").Result()
		fmt.Printf("Clearing %d clients cache entries\n", len(ks))

		for _, k := range ks {
			clientID := strings.Replace(k, "cache:clients:", "", -1)
			if Verbose {
				fmt.Printf("Clearing cached item of %s... ", clientID)
			}
			err := redisClient.Del(k).Err()
			if err != nil {
				if Verbose {
					fmt.Println("✘")
				} else {
					fmt.Printf("✘ failed to clear client cache for %s\n", clientID)
				}
			}
			if Verbose {
				fmt.Println("✔")
			}
		}
	},
}

// ✘ ✔
