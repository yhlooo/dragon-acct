package assets

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
	"github.com/yhlooo/dragon-acct/pkg/report"
)

// Options 分析选项
type Options struct {
	ShowHistory bool
}

// Analyse 分析资产数据
func Analyse(_ context.Context, assets *v1.Assets, opts Options) (report.Report, error) {
	sort.Slice(assets.Transactions, func(i, j int) bool {
		return assets.Transactions[i].Date.Before(assets.Transactions[j].Date.Time)
	})

	r := &Report{
		showHistory: opts.ShowHistory,
	}
	r.AddGoodsInfo(assets.Goods...)

	// 统计所有产品持仓情况
	allGoods := map[string]*Goods{}
	checkpointI := -1
	var checkpoint *CheckpointReport
	var checkpointGoods map[string]*Goods
	for _, t := range assets.Transactions {
		addToGoods(allGoods, t.From, true, t)
		addToGoods(allGoods, t.To, false, t)

		if checkpoint == nil || (checkpointI < len(assets.Checkpoints) && t.Date.After(checkpoint.Date.Time)) {
			// 记录当前检查点
			if checkpoint != nil {
				for _, g := range checkpointGoods {
					checkpoint.Report.goods = append(checkpoint.Report.goods, *g)
				}
				r.checkpoints = append(r.checkpoints, *checkpoint)
			}

			// 换下一个检查点
			checkpointI++
			if checkpointI < len(assets.Checkpoints) {
				checkpoint = &CheckpointReport{
					Date: assets.Checkpoints[checkpointI].Date,
				}
				for _, info := range assets.Checkpoints[checkpointI].Goods {
					detail, _ := r.GoodsInfo(info.Name)
					checkpoint.Report.AddGoodsInfo(v1.GoodsInfo{
						Name:         info.Name,
						Code:         detail.Code,
						Risk:         detail.Risk,
						Price:        info.Price,
						Base:         detail.Base,
						IgnoreReturn: detail.IgnoreReturn,
					})
				}
			} else {
				checkpoint = &CheckpointReport{
					Date: v1.Date{Time: time.Now()},
				}
				checkpoint.Report.AddGoodsInfo(assets.Goods...)
			}
			checkpointGoods = map[string]*Goods{}
		}

		addToGoods(checkpointGoods, t.From, true, t)
		addToGoods(checkpointGoods, t.To, false, t)
	}
	if checkpoint != nil {
		for _, g := range checkpointGoods {
			checkpoint.Report.goods = append(checkpoint.Report.goods, *g)
		}
		r.checkpoints = append(r.checkpoints, *checkpoint)
	}

	// 添加产品记录
	for _, g := range allGoods {
		r.goods = append(r.goods, *g)
	}

	// 补充完成
	err := r.Complete()
	if err != nil {
		return nil, fmt.Errorf("complete report error: %w", err)
	}

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
