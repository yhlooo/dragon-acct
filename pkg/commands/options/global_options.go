package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// NewDefaultGlobalOptions 返回默认全局选项
func NewDefaultGlobalOptions() GlobalOptions {
	return GlobalOptions{
		Verbosity: 0,
		Chdir:     "",
	}
}

// GlobalOptions 全局选项
type GlobalOptions struct {
	// 日志数量级别（ 0 / 1 / 2 ）
	Verbosity uint32 `json:"verbosity" yaml:"verbosity"`
	// 改变工作目录
	Chdir string `json:"chdir,omitempty" yaml:"chdir,omitempty"`
}

// Validate 校验选项是否合法
func (o *GlobalOptions) Validate() error {
	if o.Verbosity > 2 {
		return fmt.Errorf("invalid log verbosity: %d (expected: 0, 1 or 2)", o.Verbosity)
	}
	return nil
}

// AddPFlags 将选项绑定到命令行参数
func (o *GlobalOptions) AddPFlags(flags *pflag.FlagSet) {
	flags.Uint32VarP(&o.Verbosity, "verbose", "v", o.Verbosity, "Number for the log level verbosity (0, 1, or 2)")
	flags.StringVarP(&o.Chdir, "chdir", "C", o.Chdir, "Change to directory before doing anything")
}
