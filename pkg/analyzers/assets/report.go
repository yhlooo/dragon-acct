package assets

import (
	"fmt"
	"io"
	"slices"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/shopspring/decimal"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
)

// Report 报告
type Report struct {
	goodsInfos []v1.GoodsInfo
	goods      []Goods
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
	// TODO: 收益、收益率
	SortIndex uint64
}

// Complete 补充完成
func (r *Report) Complete() {
	// 产品详细信息的索引
	goodsInfos := map[string]v1.GoodsInfo{}
	goodsIndexes := map[string]int{}
	for i, info := range r.goodsInfos {
		goodsInfos[info.Name] = info
		goodsIndexes[info.Name] = i
	}

	totalValue := decimal.Zero
	for i, g := range r.goods {
		// 补充产品信息
		info, ok := goodsInfos[g.Name]
		if ok {
			r.goods[i].Code = info.Code
			r.goods[i].Price = info.Price
			r.goods[i].Risk = info.Risk
		}
		// 补充总价
		r.goods[i].Value = g.Quantity.Mul(r.goods[i].Price)
		totalValue = totalValue.Add(r.goods[i].Value)
	}

	if totalValue.Cmp(decimal.Zero) != 0 {
		for i, g := range r.goods {
			// 补充占比
			r.goods[i].Ratio = g.Value.Div(totalValue)
		}
	}

	// 排序
	sort.Slice(r.goods, func(i, j int) bool {
		aName := r.goods[i].Name
		bName := r.goods[j].Name

		// 比较在商品信息列表中的位置
		aIndex, aIndexOK := goodsIndexes[aName]
		bIndex, bIndexOK := goodsIndexes[bName]
		switch {
		case aIndexOK && !bIndexOK:
			return true
		case !aIndexOK && bIndexOK:
			return false
		case aIndex != bIndex:
			return aIndex < bIndex
		}

		// 比较持仓
		return r.goods[i].Value.Cmp(r.goods[j].Value) > 0
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
		if g.Value.Cmp(decimal.Zero) == 0 {
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
		if g.Value.Cmp(decimal.Zero) == 0 {
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
		if totalValue.Cmp(decimal.Zero) != 0 {
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

// Text 输出文本形式的报告
func (r *Report) Text(w io.Writer) error {
	r.textAllGoods(w)
	r.textHoldingGoods(w)
	r.textRisks(w)
	r.textCustodians(w)

	return nil
}

// textAllGoods 输出文本形式的关于所有产品的报告
func (r *Report) textAllGoods(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Custodian", "Code", "Risk", "Price", "Quantity", "Value"})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
	})
	for _, g := range r.AllGoods() {
		table.Append([]string{
			g.Name,
			g.Custodian,
			g.Code,
			string(g.Risk),
			g.Price.StringFixedBank(2),
			g.Quantity.StringFixedBank(2),
			g.Value.StringFixedBank(2),
		})
	}

	_, _ = fmt.Fprintln(w, "All Goods:")
	table.Render()
	_, _ = fmt.Fprintln(w)
}

// textHoldingGoods 输出文本形式的关于持仓分布的报告
func (r *Report) textHoldingGoods(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Custodian", "Value", "Ratio"})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
	})
	for _, g := range r.HoldingGoods() {
		table.Append([]string{
			g.Name,
			g.Custodian,
			g.Value.StringFixedBank(2),
			g.Ratio.Shift(2).StringFixedBank(2) + "%",
		})
	}

	_, _ = fmt.Fprintln(w, "Holding:")
	table.Render()
	_, _ = fmt.Fprintln(w)
}

// textRisks 输出文本形式的关于风险分布的报告
func (r *Report) textRisks(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Risk", "Value", "Ratio"})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
	})
	for _, g := range r.Risks() {
		table.Append([]string{
			string(g.Risk),
			g.Value.StringFixedBank(2),
			g.Ratio.Shift(2).StringFixedBank(2) + "%",
		})
	}

	_, _ = fmt.Fprintln(w, "Risks:")
	table.Render()
	_, _ = fmt.Fprintln(w)
}

// textCustodians 输出文本形式的关于托管机构分布的报告
func (r *Report) textCustodians(w io.Writer) {
	// 关键商品
	var pinned []string
	for _, info := range r.goodsInfos {
		if !info.Pinned {
			continue
		}
		pinned = append(pinned, info.Name)
	}

	// 按托管机构分组统计资产总价
	var custodians []string
	groupByCustodian := map[string]map[string]decimal.Decimal{}
	for _, g := range r.HoldingGoods() {
		name := g.Name
		if !slices.Contains(pinned, name) {
			name = "others"
		}
		if groupByCustodian[g.Custodian] == nil {
			groupByCustodian[g.Custodian] = map[string]decimal.Decimal{}
			custodians = append(custodians, g.Custodian)
		}
		groupByCustodian[g.Custodian][name] = groupByCustodian[g.Custodian][name].Add(g.Value)
		groupByCustodian[g.Custodian]["total"] = groupByCustodian[g.Custodian]["total"].Add(g.Value)
	}
	sort.Slice(custodians, func(i, j int) bool {
		return groupByCustodian[custodians[i]]["total"].Cmp(groupByCustodian[custodians[j]]["total"]) > 0
	})

	// 组装表格
	columnAlignment := []int{tablewriter.ALIGN_LEFT}
	header := []string{"Custodian"}
	var data [][]string
	for i, custodian := range custodians {
		goods := groupByCustodian[custodian]
		line := []string{custodian}

		for _, name := range pinned {
			line = append(line, goods[name].StringFixedBank(2))
			if i == 0 {
				header = append(header, name)
				columnAlignment = append(columnAlignment, tablewriter.ALIGN_RIGHT)
			}
		}

		line = append(line, goods["others"].StringFixedBank(2))
		if i == 0 {
			header = append(header, "Others")
			columnAlignment = append(columnAlignment, tablewriter.ALIGN_RIGHT)
		}

		if len(pinned) != 0 {
			line = append(line, goods["total"].StringFixedBank(2))
			if i == 0 {
				header = append(header, "Total")
				columnAlignment = append(columnAlignment, tablewriter.ALIGN_RIGHT)
			}
		}

		data = append(data, line)
	}
	table := tablewriter.NewWriter(w)
	table.SetHeader(header)
	table.SetColumnAlignment(columnAlignment)
	table.AppendBulk(data)

	_, _ = fmt.Fprintln(w, "Custodians:")
	table.Render()
	_, _ = fmt.Fprintln(w)
}
