package gormx

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type BaseModel struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ArrayAny 是一个自定义的类型，用于存储任意类型的数组, 需要对应 json 数据类型
type ArrayAny []any

func (a *ArrayAny) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return errors.New("unsupported type")
	}
}

func (a ArrayAny) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func ArrayAnyFromStrings(values []string) ArrayAny {
	var a ArrayAny
	for _, v := range values {
		a = append(a, any(v))
	}
	return a
}

func (a ArrayAny) Strings() []string {
	var values []string
	for _, v := range a {
		values = append(values, fmt.Sprintf("%v", v))
	}
	return values
}

// Any 是一个自定义的类型，用于存储任意类型的数据, 需要对应 json 数据类型
type Any map[string]any

func (a *Any) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return errors.New("unsupported type")
	}
}

func (a Any) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// StringArray 是一个自定义的类型，用于存储字符串数组, 需要对应 json 数据类型
type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return errors.New("unsupported type")
	}
}

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// StringMap 是一个自定义的类型，用于存储字符串键值对, 需要对应 json 数据类型
type StringMap map[string]any

func (s *StringMap) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return errors.New("unsupported type")
	}
}

func (s StringMap) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type IntArray []int

func (i *IntArray) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, i)
	case string:
		return json.Unmarshal([]byte(v), i)
	default:
		return errors.New("unsupported type")
	}
}

func (i IntArray) Value() (driver.Value, error) {
	return json.Marshal(i)
}

type Int64Array []int64

func (i *Int64Array) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, i)
	case string:
		return json.Unmarshal([]byte(v), i)
	default:
		return errors.New("unsupported type")
	}
}

func (i Int64Array) Value() (driver.Value, error) {
	return json.Marshal(i)
}
