package rateofreturn

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// TestXIRRCashFlowSum 测试 xirrCashFlowSum 方法
func TestXIRRCashFlowSum(t *testing.T) {
	now := time.Now()
	ret := xirrCashFlowSum([]CashFlowRecord{{
		Date:   now,
		Amount: decimal.New(-1000000, -2),
	}, {
		Date:   now.Add(43830 * time.Hour),
		Amount: decimal.New(-1000000, -2),
	}, {
		Date:   now.Add(70128 * time.Hour),
		Amount: decimal.New(501000, -2),
	}, {
		Date:   now.Add(78894 * time.Hour),
		Amount: decimal.New(1002000, -2),
	}, {
		Date:   now.Add(87660 * time.Hour),
		Amount: decimal.New(501000, -2),
	}}, decimal.New(307, -6))

	expected := decimal.New(5072167, -8)
	if !ret.Equal(expected) {
		t.Errorf("unexpected ret: %s (expected: %s)", ret, expected)
	}

	ret2 := xirrCashFlowSum([]CashFlowRecord{{
		Date:   now,
		Amount: decimal.New(-50000, 0),
	}, {
		Date:   now.Add(720 * time.Hour),
		Amount: decimal.New(-20000, 0),
	}, {
		Date:   now.Add(2160 * time.Hour),
		Amount: decimal.New(70149, 0),
	}}, decimal.New(307, -6))

	expected2 := decimal.New(14419869481, -8)
	if !ret2.Equal(expected2) {
		t.Errorf("unexpected ret: %s (expected: %s)", ret2, expected2)
	}
}

// TestXIRR 测试 XIRR 方法
func TestXIRR(t *testing.T) {
	now := time.Now()
	ret := XIRR([]CashFlowRecord{{
		Date:   now,
		Amount: decimal.New(-1000000, -2),
	}, {
		Date:   now.Add(43830 * time.Hour),
		Amount: decimal.New(-1000000, -2),
	}, {
		Date:   now.Add(70128 * time.Hour),
		Amount: decimal.New(501000, -2),
	}, {
		Date:   now.Add(78894 * time.Hour),
		Amount: decimal.New(1002000, -2),
	}, {
		Date:   now.Add(87660 * time.Hour),
		Amount: decimal.New(501000, -2),
	}})

	expected := decimal.New(307, -6)
	if !ret.Equal(expected) {
		t.Errorf("unexpected ret: %s (expected: %s)", ret, expected)
	}
}
