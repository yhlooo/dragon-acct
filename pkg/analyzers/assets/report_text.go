package assets

import (
	"fmt"
	"io"
	"slices"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/shopspring/decimal"
)

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
	table.SetHeader([]string{"Name", "Custodian", "Code", "Risk", "Price", "Quantity", "Value", "P/L", "RR", "XIRR"})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
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
			g.ProfitAndLoss.StringFixedBank(2),
			g.RateOfReturn.Mul(decimal.New(100, 0)).StringFixedBank(2) + "%",
			g.AnnualizedRateOfReturn.Mul(decimal.New(100, 0)).StringFixedBank(2) + "%",
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
		if !info.Base {
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
		return groupByCustodian[custodians[j]]["total"].LessThan(groupByCustodian[custodians[i]]["total"])
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
