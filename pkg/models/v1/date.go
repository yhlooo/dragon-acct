package v1

import (
	"encoding/json"
	"time"

	"gopkg.in/yaml.v3"
)

// Date 日期
type Date struct {
	time.Time
}

var _ yaml.Marshaler = &Date{}
var _ yaml.Unmarshaler = &Date{}
var _ json.Marshaler = &Date{}
var _ json.Unmarshaler = &Date{}

// String 返回字符串表示
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) String() string {
	return d.Format(time.DateOnly)
}

// MarshalJSON 序列化为 JSON
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON 从 JSON 反序列化
//
//goland:noinspection GoMixedReceiverTypes
func (d *Date) UnmarshalJSON(in []byte) error {
	var str string
	if err := json.Unmarshal(in, &str); err != nil {
		return err
	}
	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalYAML 序列化为 YAML
//
//goland:noinspection GoMixedReceiverTypes
func (d Date) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

// UnmarshalYAML 从 YAML 反序列化
//
//goland:noinspection GoMixedReceiverTypes
func (d *Date) UnmarshalYAML(in *yaml.Node) error {
	t, err := time.Parse(time.DateOnly, in.Value)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}
