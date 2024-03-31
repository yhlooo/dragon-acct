package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	analyzersassets "github.com/yhlooo/dragon-acct/pkg/analyzers/assets"
	"github.com/yhlooo/dragon-acct/pkg/collector"
	"github.com/yhlooo/dragon-acct/pkg/commands/options"
)

// NewRunCommandWithOptions 创建一个基于选项的 run 命令
func NewRunCommandWithOptions(opts *options.RunOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run analysis and output reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 校验选项
			if err := opts.Validate(); err != nil {
				return err
			}

			ctx := cmd.Context()

			// 获取输入
			pwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("get current workdir error: %w", err)
			}
			data, err := collector.Collect(pwd)
			if err != nil {
				return fmt.Errorf("collect error: %w", err)
			}

			// 分析
			report, err := analyzersassets.Analyse(ctx, &data.Assets)
			if err != nil {
				return err
			}

			// 输出
			switch opts.Format {
			case "text":
				w := os.Stdout
				if opts.Output != "" {
					w, err = os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE, 0o644)
					if err != nil {
						return fmt.Errorf("open output file %q error: %w", opts.Output, err)
					}
				}
				return report.Text(w)
			default:
				return fmt.Errorf("unsupported output format: %q", opts.Format)
			}
		},
	}

	// 绑定选项到命令行参数
	opts.AddPFlags(cmd.Flags())

	return cmd
}
