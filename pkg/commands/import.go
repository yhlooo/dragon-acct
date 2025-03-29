package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yhlooo/dragon-acct/pkg/commands/options"
	"github.com/yhlooo/dragon-acct/pkg/imports"
)

// NewImportCommandWithOptions 基于选项创建 import 命令
func NewImportCommandWithOptions(opts *options.ImportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [PATH]",
		Short: "Import data",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// 确定输入输出
			r := os.Stdin
			if len(args) > 0 {
				var err error
				r, err = os.Open(args[0])
				if err != nil {
					return fmt.Errorf("open %q error: %w", args[0], err)
				}
				defer func() { _ = r.Close() }()
			}
			w := os.Stdout
			if opts.Output != "" {
				var err error
				w, err = os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
				if err != nil {
					return fmt.Errorf("open %q error: %w", opts.Output, err)
				}
				defer func() { _ = w.Close() }()
			}
			switch opts.DataType {
			case "futu":
				return imports.ImportFutu(cmd.Context(), r, w)
			}
			return nil
		},
	}

	// 绑定选项到命令行参数
	opts.AddPFlags(cmd.Flags())

	return cmd

}
