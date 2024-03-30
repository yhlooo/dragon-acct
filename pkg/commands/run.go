package commands

import (
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"github.com/yhlooo/dragon-acct/pkg/commands/options"
)

// NewRunCommandWithOptions 创建一个基于选项的 run 命令
func NewRunCommandWithOptions(*options.RunOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run analysis and output reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logr.FromContextOrDiscard(cmd.Context())
			logger.Info("TODO: ...")
			return nil
		},
	}
	return cmd
}
