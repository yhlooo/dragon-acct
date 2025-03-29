package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// NewDefaultImportOptions 创建默认 import 命令选项
func NewDefaultImportOptions() ImportOptions {
	return ImportOptions{
		Output: "",
		Format: "csv",
	}
}

// ImportOptions import 命令选项
type ImportOptions struct {
	// 输出文件路径
	Output string `json:"output,omitempty" yaml:"output,omitempty"`
	// 输出格式
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
	// 导入数据类型
	DataType string `json:"dataType,omitempty" yaml:"dataType,omitempty"`
}

// Validate 校验选项是否合法
func (o *ImportOptions) Validate() error {
	switch o.Format {
	case "csv", "json", "yaml":
	default:
		return fmt.Errorf("unsupported output format: %q", o.Format)
	}
	switch o.DataType {
	case "futu":
	default:
		return fmt.Errorf("unsupported data type: %q", o.DataType)
	}
	return nil
}

// AddPFlags 将选项绑定到命令行参数
func (o *ImportOptions) AddPFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&o.DataType, "type", "t", o.DataType, `Import data type ("futu")`)
	flags.StringVarP(&o.Output, "output", "o", o.Output, "Output path")
	flags.StringVarP(&o.Format, "format", "f", o.Format, `Output format ("csv", "yaml", "json")`)
}
