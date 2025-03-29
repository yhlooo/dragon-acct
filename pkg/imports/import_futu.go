package imports

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// ImportFutu 导入富途交易记录
func ImportFutu(ctx context.Context, r io.Reader, w io.Writer) error {

	csvR := csv.NewReader(r)
	csvR.TrimLeadingSpace = true
	csvR.LazyQuotes = true
	rows, err := csvR.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv error: %w", err)
	}

	csvW := csv.NewWriter(w)
	defer csvW.Flush()

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 36 {
			return fmt.Errorf("csv row too short: %d (at least 36)", len(row))
		}
		if row[20] == "" {
			// 没有成交
			continue
		}

		side := row[0]
		name := row[2]                                                               // 产品名
		currency := row[16]                                                          // 货币
		quantity, err := decimal.NewFromString(strings.ReplaceAll(row[18], ",", "")) // 数量
		if err != nil {
			return fmt.Errorf("parse quantity %q error: %w", row[18], err)
		}
		price, err := decimal.NewFromString(strings.ReplaceAll(row[19], ",", "")) // 价格
		if err != nil {
			return fmt.Errorf("parse price %q error: %w", row[19], err)
		}
		amount, err := decimal.NewFromString(strings.ReplaceAll(row[20], ",", "")) // 金额
		if err != nil {
			return fmt.Errorf("parse amount %q error: %w", row[20], err)
		}
		date, err := time.Parse("Jan 2, 2006 15:04:05 MST", strings.ReplaceAll(row[21], "ET", "EDT")) // 日期
		if err != nil {
			return fmt.Errorf("parse date %q error: %w", row[21], err)
		}
		fees, err := decimal.NewFromString(strings.ReplaceAll(row[35], ",", "")) // 手续费
		if err != nil {
			return fmt.Errorf("parse fees %q error: %w", row[35], err)
		}

		switch side {
		case "Sell":
			_ = csvW.Write([]string{
				date.Format(time.DateOnly),
				quantity.StringFixedBank(2), name, "",
				amount.Sub(fees).StringFixedBank(2), currency, "",
				"", fmt.Sprintf("price: %s, fees: %s", price.StringFixedBank(2), fees.StringFixedBank(2)),
			})
		case "Buy":
			_ = csvW.Write([]string{
				date.Format(time.DateOnly),
				amount.Add(fees).StringFixedBank(2), currency, "",
				quantity.StringFixedBank(2), name, "",
				"", fmt.Sprintf("price: %s, fees: %s", price.StringFixedBank(2), fees.StringFixedBank(2)),
			})
		default:
			return fmt.Errorf("invalid side %q, expected 'Sell' or 'Buy'", side)
		}
	}

	return nil
}

// FutuTransaction 富途交易记录
type FutuTransaction struct {
}

// 2025-03-29 , 0.00 , 博时现金宝货币B , 汇丰银行 , 44.73 , 博时现金宝货币B , 汇丰银行 , 利息分红 , 到这天（含）累计分红
