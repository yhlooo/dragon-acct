package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// NewDefaultRunOptions 创建一个默认的 RunOptions
func NewDefaultRunOptions() RunOptions {
	return RunOptions{
		Output: "",
		Format: "text",
	}
}

// RunOptions run 命令选项
type RunOptions struct {
	// 输出文件路径
	Output string `json:"output,omitempty" yaml:"output,omitempty"`
	// 输出格式
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
}

// Validate 校验选项是否合法
func (o *RunOptions) Validate() error {
	switch o.Format {
	case "text":
	default:
		return fmt.Errorf("unsupported output format: %q", o.Format)
	}
	return nil
}

// AddPFlags 将选项绑定到命令行参数
func (o *RunOptions) AddPFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&o.Output, "output", "o", o.Output, "Output path of the report")
	flags.StringVarP(
		&o.Format, "format", "f", o.Format,
		`Output format of the report ("text", "yaml", "json", "markdown" or "html")`,
	)
}
