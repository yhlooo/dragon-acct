package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	analyzersassets "github.com/yhlooo/dragon-acct/pkg/analyzers/assets"
	analyzerincome "github.com/yhlooo/dragon-acct/pkg/analyzers/income"
	"github.com/yhlooo/dragon-acct/pkg/collector"
	"github.com/yhlooo/dragon-acct/pkg/commands/options"
	"github.com/yhlooo/dragon-acct/pkg/report"
)

// NewRunCommandWithOptions 创建一个基于选项的 run 命令
func NewRunCommandWithOptions(opts *options.RunOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [income|assets]...",
		Short: "Run analysis and output reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 校验选项
			if err := opts.Validate(); err != nil {
				return err
			}

			ctx := cmd.Context()

			targets := args
			if len(targets) == 0 {
				targets = []string{"income", "assets"}
			}

			// 获取输入
			pwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("get current workdir error: %w", err)
			}
			data, err := collector.Collect(pwd)
			if err != nil {
				return fmt.Errorf("collect error: %w", err)
			}

			for _, target := range targets {
				// 分析
				var r report.Report
				switch target {
				case "income":
					r, err = analyzerincome.Analyse(ctx, &data.Income)
				case "assets":
					r, err = analyzersassets.Analyse(ctx, &data.Assets)
				default:
					return fmt.Errorf("unsupported target: %q", target)
				}
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
					if err := r.Text(w); err != nil {
						return err
					}
				default:
					return fmt.Errorf("unsupported output format: %q", opts.Format)
				}
			}

			return nil
		},
	}

	// 绑定选项到命令行参数
	opts.AddPFlags(cmd.Flags())

	return cmd
}
