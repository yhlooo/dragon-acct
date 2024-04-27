package v1

// Root 根
type Root struct {
	// 收入
	Income Income `json:"income,omitempty" yaml:"income,omitempty"`
	// 资产
	Assets Assets `json:"assets,omitempty" yaml:"assets,omitempty"`
}
