package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cmrigney/skill-share/pkg/oci"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [skill-path] [registry/repository:tag]",
	Short: "Package and push a Claude skill to an OCI registry",
	Long: `Package a Claude skill directory into an OCI artifact and push it directly
to a container registry. The skill will not be stored locally.

Examples:
  skill-share push ./my-skill ghcr.io/username/my-skill:latest
  skill-share push ~/.claude/skills/pdf-analyzer docker.io/myuser/pdf-skill:v1.0.0
  skill-share push . ttl.sh/my-skill:1h`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		skillPath := args[0]
		ref := args[1]

		// Resolve absolute path
		absPath, err := filepath.Abs(skillPath)
		if err != nil {
			exitWithError(fmt.Errorf("failed to resolve skill path: %w", err))
		}

		// Check if path exists
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			exitWithError(fmt.Errorf("skill path does not exist: %s", absPath))
		}

		fmt.Printf("Packaging skill from: %s\n", absPath)

		// Package and push the skill
		if err := oci.PackageAndPushSkill(absPath, ref); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
