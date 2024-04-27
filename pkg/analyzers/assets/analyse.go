package assets

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
	"github.com/yhlooo/dragon-acct/pkg/report"
)

// Analyse 分析资产数据
func Analyse(_ context.Context, assets *v1.Assets) (report.Report, error) {
	// 统计所有产品持仓情况
	allGoods := map[string]*Goods{}
	for _, t := range assets.Transactions {
		addToGoods(allGoods, t.From, true, t)
		addToGoods(allGoods, t.To, false, t)
	}

	// 组装报告
	r := &Report{}
	r.AddGoodsInfo(assets.Goods...)
	for _, g := range allGoods {
		r.goods = append(r.goods, *g)
	}

	// 补充完成
	r.Complete()

	return r, nil
}

// addToGoods 添加商品交易记录
func addToGoods(allGoods map[string]*Goods, goods *v1.Goods, minus bool, t v1.Transaction) {
	if goods == nil {
		return
	}
	key := fmt.Sprintf("%s/%s", goods.Custodian, goods.Name)
	if allGoods[key] == nil {
		allGoods[key] = &Goods{
			Name:      goods.Name,
			Custodian: goods.Custodian,
			Quantity:  decimal.Zero,
		}
	}
	if minus {
		allGoods[key].Quantity = allGoods[key].Quantity.Sub(goods.Quantity)
	} else {
		allGoods[key].Quantity = allGoods[key].Quantity.Add(goods.Quantity)
	}
	allGoods[key].transactions = append(allGoods[key].transactions, t)
}
