package assets

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
	"github.com/yhlooo/dragon-acct/pkg/report"
	"github.com/yhlooo/dragon-acct/pkg/utils/rateofreturn"
)

const (
	InternalBaseGoods = "InternalBaseGoods"
)

// Report 资产报告
type Report struct {
	showHistory bool

	goodsInfos   map[string]v1.GoodsInfo
	goodsIndexes map[string]int

	goods                  []Goods
	totalValue             decimal.Decimal
	profitAndLoss          decimal.Decimal
	rateOfReturn           decimal.Decimal
	annualizedRateOfReturn decimal.Decimal

	checkpoints []CheckpointReport
}

// Goods 商品
type Goods struct {
	// 商品名
	Name string `json:"name" yaml:"name"`
	// 代号
	Code string `json:"code,omitempty" yaml:"code,omitempty"`
	// 托管机构
	Custodian string `json:"custodian,omitempty" yaml:"custodian,omitempty"`
	// 数量
	Quantity decimal.Decimal `json:"quantity" yaml:"quantity"`
	// 单价
	Price decimal.Decimal `json:"price" yaml:"price"`
	// 风险
	Risk v1.RiskLevel `json:"risk,omitempty" yaml:"risk,omitempty"`
	// 总价值
	Value decimal.Decimal `json:"value,omitempty" yaml:"value,omitempty"`
	// 占比
	Ratio decimal.Decimal `json:"ratio,omitempty" yaml:"ratio,omitempty"`
	// 损益
	ProfitAndLoss decimal.Decimal `json:"profitAndLoss,omitempty" yaml:"profitAndLoss,omitempty"`
	// 收益率
	RateOfReturn decimal.Decimal `json:"rateOfReturn,omitempty" yaml:"rateOfReturn,omitempty"`
	// 年化收益率
	AnnualizedRateOfReturn decimal.Decimal `json:"annualizedRateOfReturn,omitempty" yaml:"annualizedRateOfReturn,omitempty"`

	// 是否基础商品（货币）
	Base bool `json:"base,omitempty" yaml:"base,omitempty"`
	// 是否忽略收益
	IgnoreReturn bool `json:"ignoreReturn,omitempty" yaml:"ignoreReturn,omitempty"`

	// 关于该产品的交易
	transactions []v1.Transaction
}

var _ report.Report = &Report{}

// AddGoodsInfo 添加商品信息
func (r *Report) AddGoodsInfo(goodsInfos ...v1.GoodsInfo) {
	if len(goodsInfos) == 0 {
		return
	}

	if r.goodsInfos == nil {
		r.goodsInfos = map[string]v1.GoodsInfo{}
	}
	if r.goodsIndexes == nil {
		r.goodsIndexes = map[string]int{}
	}

	curIndex := len(r.goodsIndexes)
	for i, info := range goodsInfos {
		r.goodsInfos[info.Name] = info
		r.goodsIndexes[info.Name] = curIndex + i
	}
}

// GoodsInfo 获取商品信息
func (r *Report) GoodsInfo(name string) (v1.GoodsInfo, bool) {
	info, ok := r.goodsInfos[name]
	return info, ok
}

// Complete 补充完成
func (r *Report) Complete() error {
	r.totalValue = decimal.Zero
	for i, g := range r.goods {
		// 补充产品信息
		info, ok := r.goodsInfos[g.Name]
		if ok {
			r.goods[i].Code = info.Code
			r.goods[i].Price = info.Price
			r.goods[i].Risk = info.Risk
			r.goods[i].Base = info.Base
			r.goods[i].IgnoreReturn = info.IgnoreReturn
		}
		// 补充总价
		if !g.Quantity.IsZero() && !ok {
			return fmt.Errorf("goods %q price not found", g.Name)
		}
		r.goods[i].Value = g.Quantity.Mul(r.goods[i].Price)
		if !r.goods[i].Base || r.goods[i].Value.IsPositive() {
			r.totalValue = r.totalValue.Add(r.goods[i].Value)
		}

		// 补充损益情况
		if !r.goods[i].Base {
			totalCost, totalReturn, cashFlow, err := r.parseGoodsProfitAndLoss(&r.goods[i])
			if err != nil {
				return fmt.Errorf("complete %q with custodian %q profit and loss error: %w", g.Name, g.Custodian, err)
			}
			if totalCost.IsZero() {
				return fmt.Errorf("goods %q with custodian %q total cost is zero", g.Name, g.Custodian)
			}
			r.goods[i].ProfitAndLoss = totalReturn.Sub(totalCost)
			r.goods[i].RateOfReturn = totalReturn.Sub(totalCost).DivRound(totalCost, 6)
			r.goods[i].AnnualizedRateOfReturn = rateofreturn.XIRR(cashFlow)
		}
	}
	if err := r.completeTotalProfitAndLoss(); err != nil {
		return fmt.Errorf("complete total profit and loss error: %w", err)
	}

	if !r.totalValue.IsZero() {
		for i, g := range r.goods {
			// 补充占比
			r.goods[i].Ratio = g.Value.Div(r.totalValue)
		}
	}

	// 排序
	r.sortGoods()

	// 补充检查点
	var lastCheckpoint *CheckpointReport
	for i := range r.checkpoints {
		if err := r.checkpoints[i].Complete(lastCheckpoint); err != nil {
			return fmt.Errorf("complete checkpoint %s error: %w", r.checkpoints[i].Date, err)
		}
		lastCheckpoint = &r.checkpoints[i]
	}

	return nil
}

// completeTotalProfitAndLoss 补充总体损益情况
func (r *Report) completeTotalProfitAndLoss() error {
	var cashFlow []rateofreturn.CashFlowRecord
	totalCost := decimal.Zero
	totalReturn := decimal.Zero

	for _, goods := range r.goods {
		if goods.IgnoreReturn || goods.Base {
			continue
		}
		goodsCost, goodsReturn, goodsCashFlow, err := r.parseGoodsProfitAndLoss(&goods)
		if err != nil {
			return fmt.Errorf("get %q profit and loss error: %w", goods.Name, err)
		}
		totalCost = totalCost.Add(goodsCost)
		totalReturn = totalReturn.Add(goodsReturn)
		cashFlow = append(cashFlow, goodsCashFlow...)
	}

	r.profitAndLoss = totalReturn.Sub(totalCost)
	r.rateOfReturn = totalReturn.Sub(totalCost).DivRound(totalCost, 6)
	r.annualizedRateOfReturn = rateofreturn.XIRR(cashFlow)
	return nil
}

func (r *Report) parseGoodsProfitAndLoss(goods *Goods) (
	totalCost, totalReturn decimal.Decimal,
	cashFlow []rateofreturn.CashFlowRecord,
	err error,
) {
	for _, t := range goods.transactions {
		switch {
		case t.To == nil || t.From == nil:
			continue
		case t.To.Name == goods.Name:
			price := decimal.New(1, 0)
			if t.From.Name != InternalBaseGoods {
				info, ok := r.goodsInfos[t.From.Name]
				if !ok {
					err = fmt.Errorf("goods info %q not found", t.From.Name)
					return
				}
				price = info.Price
			}
			cashFlow = append(cashFlow, rateofreturn.CashFlowRecord{
				Date:   t.Date.Time,
				Amount: t.From.Quantity.Mul(price).Neg(),
			})
			totalCost = totalCost.Add(t.From.Quantity.Mul(price))
		case t.From.Name == goods.Name:
			price := decimal.New(1, 0)
			if t.From.Name != InternalBaseGoods {
				info, ok := r.goodsInfos[t.To.Name]
				if !ok {
					err = fmt.Errorf("goods info %q not found", t.To.Name)
					return
				}
				price = info.Price
			}
			cashFlow = append(cashFlow, rateofreturn.CashFlowRecord{
				Date:   t.Date.Time,
				Amount: t.To.Quantity.Mul(price),
			})
			totalReturn = totalReturn.Add(t.To.Quantity.Mul(price))
		}
	}
	if !goods.Value.IsZero() {
		cashFlow = append(cashFlow, rateofreturn.CashFlowRecord{
			Date:   time.Now().Round(24 * time.Hour),
			Amount: goods.Value,
		})
		totalReturn = totalReturn.Add(goods.Value)
	}
	return
}

// sortGoods 对产品进行排序
func (r *Report) sortGoods() {
	sort.Slice(r.goods, func(i, j int) bool {
		aName := r.goods[i].Name
		bName := r.goods[j].Name

		// 比较在商品信息列表中的位置
		aIndex, aIndexOK := r.goodsIndexes[aName]
		bIndex, bIndexOK := r.goodsIndexes[bName]
		switch {
		case aIndexOK && !bIndexOK:
			return true
		case !aIndexOK && bIndexOK:
			return false
		case aIndex != bIndex:
			return aIndex < bIndex
		}

		// 比较持仓
		return r.goods[j].Value.LessThan(r.goods[i].Value)
	})
}

// AllGoods 返回所有商品信息
func (r *Report) AllGoods() []Goods {
	if r.goods == nil {
		return nil
	}
	ret := make([]Goods, len(r.goods))
	copy(ret, r.goods)
	return ret
}

// HoldingGoods 返回持仓商品信息
func (r *Report) HoldingGoods() []Goods {
	var ret []Goods
	for _, g := range r.goods {
		if g.Value.IsZero() {
			continue
		}
		ret = append(ret, g)
	}
	return ret
}

// Risks 返回风险分布
func (r *Report) Risks() []Goods {
	// 统计风险分布
	risks := map[v1.RiskLevel]decimal.Decimal{}
	totalValue := decimal.Zero
	for _, g := range r.goods {
		if g.Base && g.Value.IsNegative() {
			continue
		}
		if g.Value.IsZero() {
			continue
		}
		risk := g.Risk
		if risk == "" {
			risk = "Unknown"
		}
		risks[risk] = risks[risk].Add(g.Value)
		totalValue = totalValue.Add(g.Value)
	}

	// 组装结果
	var ret []Goods
	for r, v := range risks {
		ratio := decimal.Zero
		if !totalValue.IsZero() {
			ratio = v.Div(totalValue)
		}
		ret = append(ret, Goods{
			Risk:  r,
			Value: v,
			Ratio: ratio,
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Risk < ret[j].Risk
	})
	return ret
}

// Checkpoints 返回所有检查点报告
func (r *Report) Checkpoints() []CheckpointReport {
	if len(r.checkpoints) == 0 {
		return nil
	}
	ret := make([]CheckpointReport, len(r.checkpoints))
	copy(ret, r.checkpoints)
	return ret
}

// TotalValue 返回资产总价值
func (r *Report) TotalValue() decimal.Decimal {
	return r.totalValue
}

// TotalProfitAndLoss 返回总体损益情况
func (r *Report) TotalProfitAndLoss() (profitAndLoss, rateOfReturn, annualizedRateOfReturn decimal.Decimal) {
	return r.profitAndLoss, r.rateOfReturn, r.annualizedRateOfReturn
}

// CheckpointReport 检查点报告
type CheckpointReport struct {
	// 日期
	Date v1.Date
	// 报告
	Report Report
}

// Complete 补充完成
func (r *CheckpointReport) Complete(lastCheckpoint *CheckpointReport) error {
	if lastCheckpoint == nil {
		return r.Report.Complete()
	}

	goodsMap := map[string]*Goods{}
	for i := range r.Report.goods {
		key := fmt.Sprintf("%s/%s", r.Report.goods[i].Custodian, r.Report.goods[i].Name)
		goodsMap[key] = &r.Report.goods[i]
	}

	// 添加上一期期末资产结算结果
	var extraGoods []Goods
	for _, g := range lastCheckpoint.Report.HoldingGoods() {
		key := fmt.Sprintf("%s/%s", g.Custodian, g.Name)
		if _, ok := goodsMap[key]; !ok {
			extraGoods = append(extraGoods, Goods{
				Name:      g.Name,
				Custodian: g.Custodian,
				Quantity:  decimal.Zero,
			})
			goodsMap[key] = &extraGoods[len(extraGoods)-1]
		}

		goodsMap[key].Quantity = goodsMap[key].Quantity.Add(g.Quantity)
		goodsMap[key].transactions = append(goodsMap[key].transactions, v1.Transaction{
			Date: lastCheckpoint.Date,
			From: &v1.Goods{Name: InternalBaseGoods, Quantity: g.Value},
			To:   &v1.Goods{Name: g.Name, Custodian: g.Custodian, Quantity: g.Quantity},
		})
	}
	r.Report.goods = append(r.Report.goods, extraGoods...)

	// 补全当期报告
	return r.Report.Complete()
}
