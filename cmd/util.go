package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func mustGetStringFlag(cmd *cobra.Command, name string) string {
	str, err := cmd.Flags().GetString(name)
	if err != nil {
		log.Fatalf("Error parsing flag '%s': %v", name, err)
	}

	return str
}
