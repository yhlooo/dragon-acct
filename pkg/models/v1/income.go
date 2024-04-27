package v1

import "github.com/shopspring/decimal"

// Income 收入
type Income struct {
	// 收入明细
	Details []IncomeItem `json:"details,omitempty" yaml:"details,omitempty"`
}

// IncomeItem 收入项
type IncomeItem struct {
	// 日期
	Date Date `json:"date" yaml:"date"`
	// 税前总额
	Gross decimal.Decimal `json:"gross" yaml:"gross"`
	// 保险和住房公积金
	InsuranceAndHF decimal.Decimal `json:"insuranceAndHF,omitempty" yaml:"insuranceAndHF,omitempty"`
	// 税
	Tax decimal.Decimal `json:"tax,omitempty" yaml:"tax,omitempty"`
	// 消费比例
	ConsumptionProportion decimal.Decimal `json:"consumptionProportion,omitempty" yaml:"consumptionProportion,omitempty"`
	// 标签
	Tags map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
	// 备注
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
}
