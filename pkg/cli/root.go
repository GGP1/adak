package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "palo",
	Short: "Superfree market",
	Long: `Palo lets you search for and compare products
	from the markets near your home`,
	Run: func(cmd *cobra.Command, args []string) { fmt.Println("Palo cmd") },
}

var author string

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	rootCmd.PersistentFlags().StringVar(&author, "author", "GGP", "Author name for copyright attribution")
}

// Execute executes the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
