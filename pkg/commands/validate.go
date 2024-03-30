package commands

import (
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	"github.com/yhlooo/dragon-acct/pkg/commands/options"
)

// NewValidateCommandWithOptions 创建一个基于选项的 validate 命令
func NewValidateCommandWithOptions(*options.ValidateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Check to confirm data legitimacy",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logr.FromContextOrDiscard(cmd.Context())
			logger.Info("TODO: ...")
			return nil
		},
	}
	return cmd
}
