// Package genautocomplete provides the genautocomplete command.
package genautocomplete

import (
	"github.com/Evengard/rclone/cmd"
	"github.com/spf13/cobra"
)

func init() {
	cmd.Root.AddCommand(completionDefinition)
}

var completionDefinition = &cobra.Command{
	Use:   "completion [shell]",
	Short: `Output completion script for a given shell.`,
	Long: `
Generates a shell completion script for rclone.
Run with ` + "`--help`" + ` to list the supported shells.
`,
	Annotations: map[string]string{
		"versionIntroduced": "v1.33",
	},
	Aliases: []string{"genautocomplete"},
}
