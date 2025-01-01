package collector

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v3"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
)

const (
	assetsName         = "assets"
	assetsGoods        = "assets_goods"
	assetsTransactions = "assets_transactions"
	assetsCheckpoints  = "assets_checkpoints"
	incomeName         = "income"
	incomeDetailsName  = "income_details"
)

// Collect 收集数据
func Collect(path string) (*v1.Root, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("list %q error: %w", path, err)
	}

	ret := &v1.Root{}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name())
		filePath := filepath.Join(path, f.Name())
		var err error
		switch {
		case strings.HasPrefix(f.Name(), incomeDetailsName):
			switch ext {
			case ".yaml", ".yml":
				err = loadYAML(ret, filePath, &[]v1.IncomeItem{})
			case ".csv":
				err = loadCSV(ret, filePath, &[]v1.IncomeItem{})
			}
		case strings.HasPrefix(f.Name(), incomeName):
			switch ext {
			case ".yaml", ".yml":
				err = loadYAML(ret, filePath, &v1.Income{})
			}
		case strings.HasPrefix(f.Name(), assetsCheckpoints):
			switch ext {
			case ".yaml", ".yml":
				err = loadYAML(ret, filePath, &[]v1.Checkpoint{})
			}
		case strings.HasPrefix(f.Name(), assetsTransactions):
			switch ext {
			case ".yaml", ".yml":
				err = loadYAML(ret, filePath, &[]v1.Transaction{})
			case ".csv":
				err = loadCSV(ret, filePath, &[]v1.Transaction{})
			}
		case strings.HasPrefix(f.Name(), assetsGoods):
			switch ext {
			case ".yaml", ".yml":
				err = loadYAML(ret, filePath, &[]v1.GoodsInfo{})
			case ".csv":
				err = loadCSV(ret, filePath, &[]v1.GoodsInfo{})
			}
		case strings.HasPrefix(f.Name(), assetsName):
			switch ext {
			case ".yaml", ".yml":
				err = loadYAML(ret, filePath, &v1.Assets{})
			}
		}
		if err != nil {
			return nil, fmt.Errorf("load file %q error: %w", filePath, err)
		}
	}

	return ret, nil
}

// loadYAML 加载 YAML 文件
func loadYAML(root *v1.Root, path string, into interface{}) error {
	// 读
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %q error: %w", path, err)
	}
	// 反序列化
	if err := yaml.Unmarshal(raw, into); err != nil {
		return fmt.Errorf("unmarshal file %q as yaml to %T error: %w", path, into, err)
	}
	// 合并数据
	if err := Merge(root, into); err != nil {
		return fmt.Errorf("merge file %q error: %w", path, err)
	}
	return nil
}

// loadCSV 加载 CSV 文件
func loadCSV(root *v1.Root, path string, into interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file %q error: %w", path, err)
	}
	defer func() {
		_ = f.Close()
	}()
	r := csv.NewReader(f)

	// 加载到 CSV
	switch obj := into.(type) {
	case *[]v1.IncomeItem:
		err = loadCSVToIncomeDetails(r, obj)
	case *[]v1.GoodsInfo:
		err = loadCSVToAssetsGoods(r, obj)
	case *[]v1.Transaction:
		err = loadCSVToAssetsTransactions(r, obj)
	default:
		return fmt.Errorf("can not load csv to %T", into)
	}
	if err != nil {
		return fmt.Errorf("load csv to %T error: %w", into, err)
	}

	// 合并数据
	if err := Merge(root, into); err != nil {
		return fmt.Errorf("merge file %q error: %w", path, err)
	}
	return nil
}

// loadCSVToIncomeDetails 加载 CSV 到 []v1.IncomeItem
func loadCSVToIncomeDetails(r *csv.Reader, into *[]v1.IncomeItem) error {
	rows, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv error: %w", err)
	}
	if len(rows) < 2 {
		return nil
	}

	ret := make([]v1.IncomeItem, len(rows)-1)
	for i, row := range rows[1:] {
		if len(row) != 7 {
			return fmt.Errorf("the number of columns is not as expected: %d (expected: 7)", len(row))
		}

		d, err := time.Parse(time.DateOnly, row[0])
		if err != nil {
			return fmt.Errorf("parse Date %q at line %d error: %w", row[0], i+2, err)
		}
		ret[i].Date = v1.Date{Time: d}

		ret[i].Gross, err = decimal.NewFromString(row[1])
		if err != nil {
			return fmt.Errorf("parse Gross %q at line %d error: %w", row[1], i+2, err)
		}

		ret[i].InsuranceAndHF, err = decimal.NewFromString(row[2])
		if err != nil {
			return fmt.Errorf("parse InsuranceAndHF %q at line %d error: %w", row[2], i+2, err)
		}

		ret[i].Tax, err = decimal.NewFromString(row[3])
		if err != nil {
			return fmt.Errorf("parse Tax %q at line %d error: %w", row[3], i+2, err)
		}

		ret[i].ConsumptionProportion, err = decimal.NewFromString(row[4])
		if err != nil {
			return fmt.Errorf("parse ConsumptionProportion %q at line %d error: %w", row[4], i+2, err)
		}

		if row[5] != "" {
			tags := map[string]string{}
			for _, item := range strings.Split(row[5], " ") {
				divided := strings.Split(item, ":")
				if len(divided) != 2 {
					continue
				}
				tags[divided[0]] = divided[1]
			}
			ret[i].Tags = tags
		}
		ret[i].Comment = row[6]
	}
	*into = ret

	return nil

}

// loadCSVToAssetsGoods 加载 CSV 到 []v1.GoodsInfo
func loadCSVToAssetsGoods(r *csv.Reader, into *[]v1.GoodsInfo) error {
	rows, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv error: %w", err)
	}
	if len(rows) < 2 {
		return nil
	}

	ret := make([]v1.GoodsInfo, len(rows)-1)
	for i, row := range rows[1:] {
		if len(row) != 5 {
			return fmt.Errorf("the number of columns is not as expected: %d (expected: 5)", len(row))
		}

		ret[i].Name = row[0]
		ret[i].Code = row[1]
		ret[i].Risk = v1.RiskLevel(row[2])
		ret[i].Price, err = decimal.NewFromString(row[3])
		if err != nil {
			return fmt.Errorf("parse Price %q at line %d error: %w", row[3], i+2, err)
		}
		for _, key := range strings.Split(row[4], " ") {
			switch key {
			case "Base":
				ret[i].Base = true
			case "IgnoreReturn":
				ret[i].IgnoreReturn = true
			}
		}
	}
	*into = ret
	return nil
}

// loadCSVToAssetsTransactions 加载 CSV 到 []v1.Transaction
func loadCSVToAssetsTransactions(r *csv.Reader, into *[]v1.Transaction) error {

	rows, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv error: %w", err)
	}
	if len(rows) < 2 {
		return nil
	}

	ret := make([]v1.Transaction, len(rows)-1)
	for i, row := range rows[1:] {
		if len(row) != 9 {
			return fmt.Errorf("the number of columns is not as expected: %d (expected: 9)", len(row))
		}

		d, err := time.Parse(time.DateOnly, row[0])
		if err != nil {
			return fmt.Errorf("parse Date %q at line %d error: %w", row[0], i+2, err)
		}
		ret[i].Date = v1.Date{Time: d}

		if row[2] != "" {
			quantity, err := decimal.NewFromString(row[1])
			if err != nil {
				return fmt.Errorf("parse From.Quantity %q at line %d error: %w", row[1], i+2, err)
			}
			ret[i].From = &v1.Goods{
				Quantity:  quantity,
				Name:      row[2],
				Custodian: row[3],
			}
		}
		if row[5] != "" {
			quantity, err := decimal.NewFromString(row[4])
			if err != nil {
				return fmt.Errorf("parse To.Quantity %q at line %d error: %w", row[4], i+2, err)
			}
			ret[i].To = &v1.Goods{
				Quantity:  quantity,
				Name:      row[5],
				Custodian: row[6],
			}

		}
		ret[i].Reason = row[7]
		ret[i].Comment = row[8]
	}
	*into = ret
	return nil
}
