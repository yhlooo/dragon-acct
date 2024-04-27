package report

import "io"

// Report 报告
type Report interface {
	// Text 输出文本格式的报告
	Text(w io.Writer) error
}
