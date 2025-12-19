package cmd

import (
	"github.com/cmrigney/skill-share/pkg/oci"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [registry/repository:tag] [destination-path]",
	Short: "Pull a Claude skill from an OCI registry",
	Long: `Pull a Claude skill artifact from a container registry and extract it
to a local directory.

The destination path is optional. If omitted, skills are automatically extracted
to ~/.claude/skills/<skill-name>. Specify a custom path to extract elsewhere.

Examples:
  # Pull skill (auto-extracts to ~/.claude/skills/my-skill)
  skill-share pull ghcr.io/username/my-skill:latest

  # Pull skill to custom location (e.g., for project use)
  skill-share pull ghcr.io/username/my-skill:latest ./my-skill

  # Pull to specific directory
  skill-share pull docker.io/myuser/pdf-skill:v1.0.0 ~/custom/location`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		ref := args[0]
		destPath := ""
		if len(args) > 1 {
			destPath = args[1]
		}

		// Pull and extract the skill
		if err := oci.PullSkill(ref, destPath); err != nil {
			exitWithError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
