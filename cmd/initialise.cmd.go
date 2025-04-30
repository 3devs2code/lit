package cmd

import (
	"fmt"
	"lit/internal/initialise"

	"github.com/spf13/cobra"
)

var initcmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new repo",
	Run: func(cmd *cobra.Command, args []string) {
		err := initialise.InitRepo()

		if err != nil {
			fmt.Println("Something went wrong while creating .lit folder")
		}

	},
}

func init() {
	rootCmd.AddCommand(initcmd)
}
