package cmd

import (
	"fmt"
	"lit/internal/add"

	"github.com/spf13/cobra"
)

var addcmd = &cobra.Command{
	Use:   "add",
	Short: "Add file to a repo",
	Run: func(cmd *cobra.Command, args []string) {
		message := args[0]
		fmt.Println(message)
		result, err := add.Add(message)

		if err != nil {
			fmt.Println("❌ Error adding file:", err)
			return

		}
		fmt.Println("✅ File added successfully:", result)
	},
}

func init() {
	rootCmd.AddCommand(addcmd)
}
