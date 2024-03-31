package v1

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// TestDate_MarshalYAML 测试 Date.MarshalYAML 方法
func TestDate_MarshalYAML(t *testing.T) {
	d, _ := time.Parse(time.DateOnly, "2024-12-24")
	out, err := yaml.Marshal(Date{Time: d})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(out) != "\"2024-12-24\"\n" {
		t.Errorf("unexpected result: %q (expected: \"\\\"2024-12-24\\\"\\n\")", string(out))
	}
}

// TestDate_UnmarshalYAML 测试 (*Date).UnmarshalYAML 方法
func TestDate_UnmarshalYAML(t *testing.T) {
	var d Date
	if err := yaml.Unmarshal([]byte(`"2024-12-24"`), &d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d.Year() != 2024 {
		t.Errorf("unexpected year of result: %d (expected: 2024)", d.Year())
	}
	if d.Month() != 12 {
		t.Errorf("unexpected month of result: %d (expected: 12)", d.Month())
	}
	if d.Day() != 24 {
		t.Errorf("unexpected day of result: %d (expected: 24)", d.Day())
	}
}

// TestDate_MarshalJSON 测试 Date.MarshalJSON 方法
func TestDate_MarshalJSON(t *testing.T) {
	d, _ := time.Parse(time.DateOnly, "2024-12-24")
	out, err := json.Marshal(Date{Time: d})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(out) != "\"2024-12-24\"" {
		t.Errorf("unexpected result: %q (expected: \"\\\"2024-12-24\\\"\")", string(out))
	}
}

// TestDate_UnmarshalJSON 测试 (*Date).UnmarshalJSON 方法
func TestDate_UnmarshalJSON(t *testing.T) {
	var d Date
	if err := json.Unmarshal([]byte(`"2024-12-24"`), &d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d.Year() != 2024 {
		t.Errorf("unexpected year of result: %d (expected: 2024)", d.Year())
	}
	if d.Month() != 12 {
		t.Errorf("unexpected month of result: %d (expected: 12)", d.Month())
	}
	if d.Day() != 24 {
		t.Errorf("unexpected day of result: %d (expected: 24)", d.Day())
	}
}
