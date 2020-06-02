package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		addInt(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addInt(args []string) {
	var sum int

	for _, val := range args {
		temp, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println(err)
		}
		sum = sum + temp
	}

	fmt.Printf("Sum of numbers %s is %d", args, sum)
}
