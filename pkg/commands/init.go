package commands

import (
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"github.com/yhlooo/dragon-acct/pkg/commands/options"
)

// NewInitCommandWithOptions 创建一个基于选项的 init 命令
func NewInitCommandWithOptions(*options.InitOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes project in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logr.FromContextOrDiscard(cmd.Context())
			logger.Info("TODO: ...")
			return nil
		},
	}
	return cmd
}
