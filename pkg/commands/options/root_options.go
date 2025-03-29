package options

// NewDefaultOptions 创建一个默认运行选项
func NewDefaultOptions() Options {
	return Options{
		Global:   NewDefaultGlobalOptions(),
		Init:     NewDefaultInitOptions(),
		Run:      NewDefaultRunOptions(),
		Validate: NewDefaultValidateOptions(),
		Import:   NewDefaultImportOptions(),
	}
}

// Options 运行选项
type Options struct {
	// 全局选项
	Global GlobalOptions `json:"global,omitempty" yaml:"global,omitempty"`
	// init 命令选项
	Init InitOptions `json:"init,omitempty" yaml:"init,omitempty"`
	// run 命令选项
	Run RunOptions `json:"run,omitempty" yaml:"run,omitempty"`
	// validate 命令选项
	Validate ValidateOptions `json:"validate,omitempty" yaml:"validate"`
	// import 命令选项
	Import ImportOptions `json:"import,omitempty" yaml:"import,omitempty"`
}
