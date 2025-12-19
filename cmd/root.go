package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skill-share",
	Short: "Package and share Claude skills as OCI artifacts",
	Long: `skill-share is a CLI tool for packaging personal Claude skills into OCI artifacts
and pushing them to container registries for easy sharing and distribution.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
