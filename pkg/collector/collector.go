package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	v1 "github.com/yhlooo/dragon-acct/pkg/models/v1"
)

const (
	assetsName = "assets"
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
		switch {
		case strings.HasPrefix(f.Name(), assetsName+"."):
			switch ext {
			case ".yaml", ".yml":
				raw, err := os.ReadFile(filepath.Join(path, f.Name()))
				if err != nil {
					return nil, fmt.Errorf("read file %q error: %w", f.Name(), err)
				}
				var assets v1.Assets
				if err := yaml.Unmarshal(raw, &assets); err != nil {
					return nil, fmt.Errorf("unmarshal assets file %q error: %w", f.Name(), err)
				}
				ret.Assets = assets
			}
		}
	}

	return ret, nil
}
