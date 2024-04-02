package collector

import (
	"fmt"
	"sort"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
)

// Merge 将 data 合并到 root
func Merge(root *v1.Root, data interface{}) error {
	switch d := data.(type) {
	case *v1.Assets:
		return mergeAssets(root, d)
	case *[]v1.GoodsInfo:
		return mergeAssetsGoods(root, *d)
	case *[]v1.Transaction:
		return mergeAssetsTransactions(root, *d)
	case []v1.GoodsInfo:
		return mergeAssetsGoods(root, d)
	case []v1.Transaction:
		return mergeAssetsTransactions(root, d)
	default:
		return fmt.Errorf("can not merge %T to *v1.Root", data)
	}
}

// mergeAssets 将 data 合并到 root.Assets
func mergeAssets(root *v1.Root, data *v1.Assets) error {
	if err := mergeAssetsGoods(root, data.Goods); err != nil {
		return err
	}
	if err := mergeAssetsTransactions(root, data.Transactions); err != nil {
		return err
	}
	return nil
}

// mergeAssetsGoods 将 data 合并到 root.Assets.Goods
func mergeAssetsGoods(root *v1.Root, data []v1.GoodsInfo) error {
	allGoods := make(map[string]v1.GoodsInfo, len(root.Assets.Goods)+len(data))
	for _, g := range root.Assets.Goods {
		allGoods[g.Name] = g
	}
	// 检查是否重复
	for _, g := range data {
		if _, ok := allGoods[g.Name]; ok {
			return fmt.Errorf("duplicate goods: %q", g.Name)
		}
		allGoods[g.Name] = g
	}
	// 追加
	root.Assets.Goods = append(root.Assets.Goods, data...)
	return nil
}

// mergeAssetsTransactions 将 data 合并到 root.Assets.Transactions
func mergeAssetsTransactions(root *v1.Root, data []v1.Transaction) error {
	// 追加
	root.Assets.Transactions = append(root.Assets.Transactions, data...)
	// 排序
	sort.Slice(root.Assets.Transactions, func(i, j int) bool {
		return root.Assets.Transactions[i].Date.Before(root.Assets.Transactions[j].Date.Time)
	})
	return nil
}
