package rateofreturn

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	xirrAccuracyExp   = -6 // 1E-6
	xirrMaxIterations = 100000
)

// CashFlowRecord 现金流记录
type CashFlowRecord struct {
	// 日期
	Date time.Time
	// 金额
	Amount decimal.Decimal
}

// XIRR 使用 XIRR 算法计算一组现金流的收益率
func XIRR(cashFlow []CashFlowRecord) decimal.Decimal {
	minInterval := decimal.New(1, xirrAccuracyExp)
	minRate := decimal.New(-1, 0).Add(minInterval)
	maxRate := decimal.New(1000, 0)

	// 确定迭代上界
	upper := minRate
	step := decimal.New(1, 0)
	for xirrCashFlowSum(cashFlow, upper).IsPositive() {
		step = step.Mul(decimal.New(2, 0))
		upper = upper.Add(step)
		if maxRate.LessThan(upper) {
			// 超出上限
			return maxRate
		}
	}
	if upper.LessThanOrEqual(minRate) {
		// 无解
		return minRate
	}

	// 确定迭代下界
	lower := upper.Sub(step)

	// 迭代逼近
	totalIterations := 0
	for minInterval.LessThan(upper.Sub(lower)) {
		mid := upper.Add(lower).DivRound(decimal.New(2, 0), -xirrAccuracyExp)
		if xirrCashFlowSum(cashFlow, mid).IsPositive() {
			lower = mid
		} else {
			upper = mid
		}

		totalIterations++
		if totalIterations > xirrMaxIterations {
			panic(fmt.Errorf("max iterations. cash flow: %v", cashFlow))
		}
	}

	// 达到预期精度
	return lower
}

// xirrCashFlowSum 使用猜测的利率计算 XIRR 现金流总和
func xirrCashFlowSum(cashFlow []CashFlowRecord, rate decimal.Decimal) decimal.Decimal {
	if len(cashFlow) == 0 {
		return decimal.Zero
	}

	firstDate := cashFlow[0].Date
	for _, record := range cashFlow[1:] {
		if record.Date.Before(firstDate) {
			firstDate = record.Date
		}
	}

	ret := decimal.Zero
	for _, record := range cashFlow {
		t := decimal.New(record.Date.Sub(firstDate).Milliseconds(), 0).
			DivRound(decimal.New(1000*3600*24*36525, -2), 8) // 间隔年数，以一年 365.25 天计算
		ret = ret.Add(record.Amount.DivRound(rate.Add(decimal.New(1, 0)).Pow(t), 8))
	}
	return ret
}
