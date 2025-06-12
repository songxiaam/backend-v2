package utility

import (
	"fmt"
	"reflect"
)

// Init all the utility package which need init at startup
func Init() (err error) {
	err = initSequence()
	if err != nil {
		return
	}

	return
}

func ConvertToInterfaceSlice[T any](slice []T) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// 泛型函数：将结构体切片转为 map[K]T
func SliceToMapByField[K comparable, T any](slice []T, keyField string) (map[K]T, error) {
	result := make(map[K]T)

	for _, item := range slice {
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected struct, got %s", val.Kind())
		}

		field := val.FieldByName(keyField)
		if !field.IsValid() {
			return nil, fmt.Errorf("field %s not found in struct", keyField)
		}
		if !field.CanInterface() {
			return nil, fmt.Errorf("field %s cannot be accessed", keyField)
		}

		key, ok := field.Interface().(K)
		if !ok {
			return nil, fmt.Errorf("field %s cannot be converted to key type", keyField)
		}

		result[key] = item
	}

	return result, nil
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
