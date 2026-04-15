package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONTime 自定义时间类型
type JSONTime struct {
	time.Time
}

const TimeFormat = "2006-01-02 15:04:05"

// MarshalJSON 序列化为 JSON
func (t *JSONTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf(`"%s"`, t.Format(TimeFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON 反序列化 JSON
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	// 处理 null 值
	if len(data) == 4 && string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}

	// 去掉引号
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// 处理空字符串
	if str == "" {
		t.Time = time.Time{}
		return nil
	}

	// 解析时间
	parsed, err := time.ParseInLocation(TimeFormat, str, time.Local)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// Value 数据库写入
func (t *JSONTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 数据库读取
func (t *JSONTime) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into JSONTime", value)
	}
}

// String 格式化时间字符串
func (t *JSONTime) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TimeFormat)
}

// Now 返回当前时间
func Now() JSONTime {
	return JSONTime{Time: time.Now().Truncate(time.Second)}
}

// NewJSONTime 创建 JSONTime
func NewJSONTime(t time.Time) JSONTime {
	return JSONTime{Time: t}
}

// ToJSONTime 转换为 JSONTime
func ToJSONTime(t *time.Time) *JSONTime {
	if t == nil {
		return nil
	}
	return &JSONTime{Time: *t}
}

// FromJSONTime 转换为 time.Time
func FromJSONTime(t *JSONTime) *time.Time {
	if t == nil || t.IsZero() {
		return nil
	}
	return &t.Time
}
