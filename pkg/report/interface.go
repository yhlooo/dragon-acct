package report

import "io"

// TextOptions 以文本格式输出的选项
type TextOptions struct {
	WithColor bool
}

// Report 报告
type Report interface {
	// Text 输出文本格式的报告
	Text(w io.Writer, opts TextOptions) error
}
