package income

import (
	"context"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
	"github.com/yhlooo/dragon-acct/pkg/report"
)

// Analyse 分析收入数据
func Analyse(_ context.Context, income *v1.Income) (report.Report, error) {
	details := make([]IncomeItem, len(income.Details))
	for i, item := range income.Details {
		details[i].IncomeItem = item
	}
	r := &Report{details: details}
	r.Complete()
	return r, nil
}
