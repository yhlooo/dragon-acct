package income

import (
	"sort"

	"github.com/shopspring/decimal"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
	"github.com/yhlooo/dragon-acct/pkg/report"
)

// Report 收入报告
type Report struct {
	details []IncomeItem
}

var _ report.Report = &Report{}

// IncomeItem 收入项
//
//goland:noinspection GoNameStartsWithPackageName
type IncomeItem struct {
	v1.IncomeItem

	TagValue string
	// 到手收入
	TakeHome decimal.Decimal
	// 用于消费的数量
	Consumption decimal.Decimal
}

// Complete 补充完成
func (r *Report) Complete() {
	for i, item := range r.details {
		r.details[i].TakeHome = item.Gross.Sub(item.InsuranceAndHF).Sub(item.Tax)
		r.details[i].Consumption = r.details[i].TakeHome.Mul(item.ConsumptionProportion)
	}
}

// GroupByTags 返回按标签聚合的收入数据
func (r *Report) GroupByTags() map[string][]IncomeItem {
	var tagsMap map[string]map[string]IncomeItem

	for _, item := range r.details {
		if len(item.Tags) == 0 {
			// 没有 tag
			continue
		}
		if tagsMap == nil {
			tagsMap = make(map[string]map[string]IncomeItem)
		}

		for k, v := range item.Tags {
			if tagsMap[k] == nil {
				tagsMap[k] = make(map[string]IncomeItem)
			}

			totalTakeHome := tagsMap[k][v].TakeHome.Add(item.TakeHome)
			totalConsumption := tagsMap[k][v].Consumption.Add(item.Consumption)
			tagsMap[k][v] = IncomeItem{
				IncomeItem: v1.IncomeItem{
					Gross:                 tagsMap[k][v].Gross.Add(item.Gross),
					InsuranceAndHF:        tagsMap[k][v].InsuranceAndHF.Add(item.InsuranceAndHF),
					Tax:                   tagsMap[k][v].Tax.Add(item.Tax),
					ConsumptionProportion: totalConsumption.Div(totalTakeHome),
				},
				TakeHome:    totalTakeHome,
				Consumption: totalConsumption,
			}
		}
	}

	ret := map[string][]IncomeItem{}
	for tag, values := range tagsMap {
		ret[tag] = make([]IncomeItem, 0, len(values))
		for tagValue, data := range values {
			data.TagValue = tagValue
			ret[tag] = append(ret[tag], data)
		}
		sort.Slice(ret[tag], func(i, j int) bool {
			return ret[tag][i].TagValue < ret[tag][j].TagValue
		})
	}

	return ret
}
