package v1

import (
	"github.com/shopspring/decimal"
)

// Assets 资产
type Assets struct {
	// 商品信息
	Goods []GoodsInfo `json:"goods,omitempty" yaml:"goods,omitempty"`
	// 交易记录
	Transactions []Transaction `json:"transactions,omitempty" yaml:"transactions,omitempty"`
}

// Transaction 交易
type Transaction struct {
	// 交易日期
	Date Date `json:"data" yaml:"date"`
	// 源商品
	From *Goods `json:"from,omitempty" yaml:"from,flow,omitempty"`
	// 目标商品
	To *Goods `json:"to,omitempty" yaml:"to,flow,omitempty"`
	// 交易原因
	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`
	// 备注
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
}

// Goods 商品（交易物）
type Goods struct {
	// 数量
	Quantity decimal.Decimal `json:"quantity,omitempty" yaml:"quantity,omitempty"`
	// 商品名
	Name string `json:"name" yaml:"name"`
	// 托管机构
	Custodian string `json:"custodian,omitempty" yaml:"custodian,omitempty"`
}

// GoodsInfo 商品信息
type GoodsInfo struct {
	// 商品名
	Name string `json:"name" yaml:"name"`
	// 代号
	Code string `json:"code,omitempty" yaml:"code,omitempty"`
	// 风险
	Risk RiskLevel `json:"risk,omitempty" yaml:"risk,omitempty"`
	// 单价
	Price decimal.Decimal `json:"price" yaml:"price"`
	// 图钉钉住
	Pinned bool `json:"pinned,omitempty" yaml:"pin,omitempty"`
}

// RiskLevel 风险级别
type RiskLevel string

// RiskLevel 的可选值
const (
	Risk0 = "R0"
	Risk1 = "R1"
	Risk2 = "R2"
	Risk3 = "R3"
	Risk4 = "R4"
	Risk5 = "R5"
)
