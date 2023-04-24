package tls

import (
	"reflect"
)

// ParserColumns parser obj the tag columns names
func ParserColumns(obj any, tag string, keyword bool) []string {
	result := make([]string, 0)

	valueOf := reflect.ValueOf(obj)
	typeOf := reflect.TypeOf(obj)

	if typeOf.Kind() == reflect.Ptr {
		valueOf = reflect.ValueOf(obj).Elem()
		typeOf = typeOf.Elem()
	}

	numField := valueOf.NumField()
	for i := 0; i < numField; i++ {
		tag := typeOf.Field(i).Tag.Get(tag)
		if len(tag) > 0 && tag != "-" {
			if keyword {
				result = append(result, KeywordTo(tag))
				continue
			}
			result = append(result, tag)
		}
	}
	return result
}

// ParserTagToMap read obj the tag generate map[string]interface{}
func ParserTagToMap(obj any, tag string) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if obj != nil {
		valueOf := reflect.ValueOf(obj)
		typeOf := reflect.TypeOf(obj)
		if typeOf.Kind() == reflect.Ptr {
			valueOf = reflect.ValueOf(obj).Elem()
			typeOf = reflect.TypeOf(obj).Elem()
		}

		numField := valueOf.NumField()
		for i := 0; i < numField; i++ {
			tag := typeOf.Field(i).Tag.Get(tag)
			if len(tag) > 0 && tag != "-" {
				params[tag] = nil
				field := valueOf.Field(i)
				if field.Kind() == reflect.Ptr {
					field = field.Elem()
				}
				if field.IsValid() {
					params[tag] = field.Interface()
				}
			}
		}
	}
	return params, nil
}

// ParserClause
func ParserClause(paramMap map[string]interface{}, keyword bool, filterateColumns ...string) ([]string, []interface{}) {
	resultMap := make(map[string]interface{})
	for key := range paramMap {
		having := false
		for _, columns := range filterateColumns {
			if key == columns {
				having = true
				break
			}
		}
		if !having {
			resultMap[key] = paramMap[key]
		}
	}

	resultKeys := make([]string, 0)
	resultVals := make([]interface{}, 0)

	for key, val := range resultMap {
		if keyword {
			resultKeys = append(resultKeys, KeywordTo(key))
		} else {
			resultKeys = append(resultKeys, key)
		}
		resultVals = append(resultVals, val)
	}

	return resultKeys, resultVals
}

// KeywordTo
func KeywordTo(key string) string {
	if key != "" {
		if len(key) > 1 && key[0] != '`' {
			return "`" + key + "`"
		}
	}
	return key
}
