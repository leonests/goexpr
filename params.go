package goexpr

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func buildPathFromRight(right interface{}, path []string) (res []string, err error) {
	if right == nil {
		return path, nil
	}
	defer func() {
		if r := recover(); r != nil {
			expr := strings.Join(path, ".")
			err = fmt.Errorf("failed to build parameter path %s: %v", expr, r)
			res = nil
		}
	}()

	val := reflect.ValueOf(right)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for _, r := range right.([]interface{}) {
			rv := reflect.ValueOf(r)
			if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
				path = append(path, r.([]string)...)
			} else if rv.Kind() == reflect.Float64 {
				path = append(path, strconv.Itoa(int(r.(float64))))
			} else if rv.Kind() == reflect.String {
				path = append(path, r.(string))
			} else {
				return nil, fmt.Errorf("invalid right value type %s", rv.Kind().String())
			}
		}
	case reflect.String:
		path = append(path, right.(string))
	case reflect.Float64:
		path = append(path, strconv.Itoa(int(right.(float64))))
	default:
		return nil, fmt.Errorf("invalid right value type %s", val.Kind().String())
	}

	return path, nil
}

func extractValueFromParams(params map[string]interface{}, path []string) (res interface{}, err error) {
	expr := strings.Join(path, ".")

	if len(path) == 0 {
		return nil, fmt.Errorf("invalid selector path " + expr)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to access %s: %v", expr, r)
			res = nil
		}
	}()

	value, ok := params[path[0]]
	if !ok {
		return nil, fmt.Errorf("no parameter " + path[0] + "found")
	}

	for i := 1; i < len(path); i++ {
		val := reflect.ValueOf(value)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		switch val.Kind() {
		case reflect.Struct:
			v := val.FieldByName(path[i])
			if v != (reflect.Value{}) {
				value = v.Interface()
				continue
			}
		case reflect.Map:
			v := val.MapIndex(reflect.ValueOf(path[i]))
			if v != (reflect.Value{}) {
				value = v.Interface()
				continue
			}
		case reflect.Slice, reflect.Array:
			idx, err := strconv.Atoi(path[i])
			if err != nil {
				return nil, fmt.Errorf("slice index must be int, not '%v'", path[i])
			}
			v := val.Index(idx)
			if v != (reflect.Value{}) {
				value = v.Interface()
				continue
			}
		case reflect.String:
			idx, err := strconv.Atoi(path[i])
			if err != nil {
				return nil, fmt.Errorf("string slice index must be int, not '%v'", path[i])
			}
			v := val.String()
			if v != "" {
				value = []rune(v)[idx]
				continue
			}
		default:
			return nil, fmt.Errorf("invalid type %v for selector", val.Kind().String())
		}
	}
	return convert2Float64(value), nil
}

func convert2Float64(value interface{}) interface{} {
	switch val := value.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	// case int32: //rune is int32,
	// 	return float64(val)
	case int64:
		return float64(val)
	case float32:
		return float64(val)
	case uint:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	default:
		return value
	}
}
