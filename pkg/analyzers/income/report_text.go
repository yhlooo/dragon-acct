package income

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/shopspring/decimal"
)

// Text 输出文本形式的报告
func (r *Report) Text(w io.Writer) error {
	r.textDetails(w)
	r.textGroupByTags(w)

	return nil
}

// textDetails 输出文本形式的收入明细报告
func (r *Report) textDetails(w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"Date",
		"Gross", "Insurance & HF", "Tax", "Take Home",
		"%Consumption", "Consumption",
		"Tags", "Comment",
	})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
	})
	for _, g := range r.details {
		tags := make([]string, 0, len(g.Tags))
		for k, v := range g.Tags {
			tags = append(tags, fmt.Sprintf("%s:%s", k, v))
		}
		sort.Strings(tags)

		table.Append([]string{
			g.Date.String(),
			g.Gross.StringFixedBank(2),
			g.InsuranceAndHF.StringFixedBank(2),
			g.Tax.StringFixedBank(2),
			g.TakeHome.StringFixedBank(2),
			g.ConsumptionProportion.Mul(decimal.New(100, 0)).StringFixedBank(2) + "%",
			g.Consumption.StringFixedBank(2),
			strings.Join(tags, " "),
			g.Comment,
		})
	}

	_, _ = fmt.Fprintln(w, "Details:")
	table.Render()
	_, _ = fmt.Fprintln(w)
}

// textGroupByTags 输出文本格式的按标签聚合的收入报告
func (r *Report) textGroupByTags(w io.Writer) {
	data := r.GroupByTags()
	if data == nil {
		return
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		values := make([]string, 0, len(data[k]))
		for v := range data[k] {
			values = append(values, v)
		}

		table := tablewriter.NewWriter(w)
		table.SetHeader([]string{
			k,
			"Gross", "Insurance & HF", "Tax", "Take Home",
			"%Consumption", "Consumption",
		})
		table.SetColumnAlignment([]int{
			tablewriter.ALIGN_LEFT,
			tablewriter.ALIGN_RIGHT,
			tablewriter.ALIGN_RIGHT,
			tablewriter.ALIGN_RIGHT,
			tablewriter.ALIGN_RIGHT,
			tablewriter.ALIGN_RIGHT,
			tablewriter.ALIGN_RIGHT,
		})
		for _, v := range values {
			item := data[k][v]
			table.Append([]string{
				v,
				item.Gross.StringFixedBank(2),
				item.InsuranceAndHF.StringFixedBank(2),
				item.Tax.StringFixedBank(2),
				item.TakeHome.StringFixedBank(2),
				item.ConsumptionProportion.Mul(decimal.New(100, 0)).StringFixedBank(2) + "%",
				item.Consumption.StringFixedBank(2),
			})
		}

		_, _ = fmt.Fprintf(w, "Group by %s:\n", k)
		table.Render()
		_, _ = fmt.Fprintln(w)
	}
}
