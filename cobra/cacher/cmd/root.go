package cmd

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/spf13/cobra"

	"github.com/mycodesmells/golang-examples/cobra/cacher/rdb"
)

var DBAddr, DBPassword string
var Verbose bool

var redisClient *redis.Client

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVar(&DBAddr, "addr", "localhost:6379", "address of Redis database")
	RootCmd.PersistentFlags().StringVar(&DBPassword, "pass", "", "password for Redis database")

	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(clearCmd)
}

var RootCmd = &cobra.Command{
	Use:              "cacher",
	PersistentPreRun: redisConnect,
}

func redisConnect(cmd *cobra.Command, args []string) {
	client, err := rdb.Connect(DBAddr, DBPassword)
	if err != nil {
		fmt.Println("✘ Failed to connect to Redis, configuration is not correct")
		os.Exit(1)
	}
	redisClient = client
}
