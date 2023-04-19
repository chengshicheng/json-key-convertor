package jsonkeyconvertor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/iancoleman/strcase"
)

const (
	Upper = "upper" // CityId -> CITYID
	Lower = "lower" // CityId -> cityid
	Snake = "snake" // CityId -> city_id
	Camel = "camel" // city_id -> CityId
)

var convertFuncsMap sync.Map

func init() {
	convertFuncsMap.Store(Upper, toUpperFunc)
	convertFuncsMap.Store(Lower, toLowerFunc)
	convertFuncsMap.Store(Snake, toSnakeFunc)
	convertFuncsMap.Store(Camel, toCamelFunc)
}

// RegisterConvertFunc register convert function
func RegisterConvertFunc(name string, f func(string) string) {
	if f == nil {
		panic("cannot register function as nil")
	}
	_, ok := convertFuncsMap.Load(name)
	if ok {
		panic(fmt.Sprintf("convert function already registered, convertName:%s", name))
	}
	convertFuncsMap.Store(name, f)
}

// GetConvertFunc get convert function
func GetConvertFunc(name string) (f func(string) string, ok bool) {
	value, ok := convertFuncsMap.Load(name)
	if ok {
		return value.(func(string) string), ok
	}
	return nil, false
}

// ConvertKey convert json key , default convert function is Upper
func ConvertKey(obj []byte, name string) ([]byte, error) {
	if !json.Valid(obj) {
		return nil, errors.New("invaild json")
	}
	if name == "" {
		name = Upper
	}
	f, ok := GetConvertFunc(name)
	if !ok {
		return nil, errors.New(fmt.Sprintf("convert function not registered, convertName:%s", name))
	}
	return convertKeyWithFunc(obj, f)
}

func convertKeyWithFunc(obj []byte, f func(string) string) ([]byte, error) {
	if bytes.HasPrefix(obj, []byte("[")) && bytes.HasSuffix(obj, []byte("]")) && bytes.Contains(obj, []byte("{")) {
		// check if it is a json array
		var listobj []json.RawMessage
		err := json.Unmarshal(obj, &listobj)
		if err != nil {
			return nil, err
		}
		for i, v := range listobj {
			listobj[i], err = convertKeyWithFunc(v, f)
			if err != nil {
				return nil, err
			}
		}
		return json.Marshal(listobj)
	} else if bytes.HasPrefix(obj, []byte("{")) && bytes.HasSuffix(obj, []byte("}")) {
		// check if it is a json object
		var m map[string]json.RawMessage
		err := json.Unmarshal(obj, &m)
		if err != nil {
			return nil, err
		}
		for k, v := range m {
			newK := f(k)
			delete(m, k)
			m[newK], err = convertKeyWithFunc(v, f)
			if err != nil {
				return nil, err
			}
		}
		return json.Marshal(m)
	} else {
		// is only json value, return directly
		return obj, nil
	}
}

// toUpperFunc convert json key to upper
func toUpperFunc(jsonStr string) string {
	return strings.ToUpper(jsonStr)
}

// toLowerFunc convert json key to lower
func toLowerFunc(jsonStr string) string {
	return strings.ToLower(jsonStr)
}

// toSnakeFunc convert json key to snake case
func toSnakeFunc(s string) string {
	return strcase.ToSnake(s)
}

// toCamelFunc convert json key to camel case
func toCamelFunc(s string) string {
	return strcase.ToCamel(s)
}
