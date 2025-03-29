package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yhlooo/dragon-acct/pkg/commands/options"
	cmdutil "github.com/yhlooo/dragon-acct/pkg/utils/cmd"
)

// NewDragonCommandWithOptions 创建一个基于选项的 dragon 命令
func NewDragonCommandWithOptions(opts options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "dragon",
		Short:        "A dragon interested in counting gold.",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 校验全局选项
			if err := opts.Global.Validate(); err != nil {
				return err
			}
			// 设置日志
			logger := cmdutil.SetLogger(cmd, opts.Global.Verbosity)
			// 设置工作目录
			if err := cmdutil.ChangeWorkingDirectory(cmd, opts.Global.Chdir); err != nil {
				return err
			}

			logger.V(1).Info(fmt.Sprintf("command: %q, args: %#v, options: %#v", cmd.Name(), args, opts))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// 绑定选项到命令行参数
	opts.Global.AddPFlags(cmd.PersistentFlags())

	// 添加子命令
	cmd.AddCommand(
		NewInitCommandWithOptions(&opts.Init),
		NewValidateCommandWithOptions(&opts.Validate),
		NewRunCommandWithOptions(&opts.Run),
		NewImportCommandWithOptions(&opts.Import),
	)

	return cmd
}

// NewDragonCommand 使用默认选项创建一个 dragon 命令
func NewDragonCommand() *cobra.Command {
	return NewDragonCommandWithOptions(options.NewDefaultOptions())
}
