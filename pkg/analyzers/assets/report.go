package assets

import (
	"sort"
	"time"

	"github.com/shopspring/decimal"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
	"github.com/yhlooo/dragon-acct/pkg/utils/rateofreturn"
)

// Report 报告
type Report struct {
	goodsInfos   map[string]v1.GoodsInfo
	goodsIndexes map[string]int

	//goodsInfos []v1.GoodsInfo
	goods []Goods
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

	// 关于该产品的交易
	transactions []v1.Transaction
}

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

// Complete 补充完成
func (r *Report) Complete() {
	totalValue := decimal.Zero
	for i, g := range r.goods {
		// 补充产品信息
		info, ok := r.goodsInfos[g.Name]
		if ok {
			r.goods[i].Code = info.Code
			r.goods[i].Price = info.Price
			r.goods[i].Risk = info.Risk
		}
		// 补充总价
		r.goods[i].Value = g.Quantity.Mul(r.goods[i].Price)
		totalValue = totalValue.Add(r.goods[i].Value)
		// 补充损益情况
		r.completeGoodsProfitAndLoss(&r.goods[i])
	}

	if !totalValue.IsZero() {
		for i, g := range r.goods {
			// 补充占比
			r.goods[i].Ratio = g.Value.Div(totalValue)
		}
	}

	// 排序
	r.sortGoods()
}

// completeGoodsProfitAndLoss 补充损益情况
func (r *Report) completeGoodsProfitAndLoss(goods *Goods) {
	var cashFlow []rateofreturn.CashFlowRecord
	totalCost := decimal.Zero
	totalReturn := decimal.Zero

	for _, t := range goods.transactions {
		switch {
		case t.To == nil || t.From == nil:
			continue
		case t.To.Name == goods.Name:
			price := decimal.New(1, 0)
			if info, ok := r.goodsInfos[t.From.Name]; ok {
				price = info.Price
			}
			cashFlow = append(cashFlow, rateofreturn.CashFlowRecord{
				Date:   t.Date.Time,
				Amount: t.From.Quantity.Mul(price).Neg(),
			})
			totalCost = totalCost.Add(t.From.Quantity.Mul(price))
		case t.From.Name == goods.Name:
			price := decimal.New(1, 0)
			if info, ok := r.goodsInfos[t.To.Name]; ok {
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

	goods.ProfitAndLoss = totalReturn.Sub(totalCost)
	goods.RateOfReturn = totalReturn.Sub(totalCost).DivRound(totalCost, 6)
	goods.AnnualizedRateOfReturn = rateofreturn.XIRR(cashFlow)
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
