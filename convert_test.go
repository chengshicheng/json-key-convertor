package jsonkeyconvertor

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertKey_InvalidInput(t *testing.T) {
	input := []byte(`{"foo_bar": "baz",`)
	_, err := ConvertKey(input, Upper)
	a := assert.New(t)
	a.EqualError(err, "invaild json")
}

func TestConvertFunc(t *testing.T) {
	// init test data
	type Input struct {
		InputStr    string
		convertName string
	}

	data := []struct {
		input  Input
		output string
	}{
		{input: Input{InputStr: "name", convertName: Upper}, output: "NAME"},
		{input: Input{InputStr: "Name", convertName: Upper}, output: "NAME"},
		{input: Input{InputStr: "name_", convertName: Upper}, output: "NAME_"},

		{input: Input{InputStr: "NAME", convertName: Lower}, output: "name"},
		{input: Input{InputStr: "Name", convertName: Lower}, output: "name"},
		{input: Input{InputStr: "NAME_", convertName: Lower}, output: "name_"},

		{input: Input{InputStr: "cityName", convertName: Snake}, output: "city_name"},
		{input: Input{InputStr: "CityName", convertName: Snake}, output: "city_name"},
		{input: Input{InputStr: "_CityName", convertName: Snake}, output: "_city_name"},
		{input: Input{InputStr: "CITY_Name", convertName: Snake}, output: "city_name"},
		{input: Input{InputStr: "CITY_ID", convertName: Snake}, output: "city_id"},

		{input: Input{InputStr: "city", convertName: Camel}, output: "City"},
		{input: Input{InputStr: "cityname", convertName: Camel}, output: "Cityname"},
		{input: Input{InputStr: "city_name", convertName: Camel}, output: "CityName"},
	}

	for _, tt := range data {
		t.Run(tt.input.convertName, func(t *testing.T) {
			f, ok := GetConvertFunc(tt.input.convertName)
			if !ok {
				t.Errorf("%s func not register", tt.input.convertName)
			}
			result := f(tt.input.InputStr)
			// 判断结果是否符合期望
			if result != tt.output {
				t.Errorf("expected %s, got %s", tt.output, result)
			}
		})
	}
}

func TestGetConvertFunc(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		args   args
		wantF  func(string) string
		wantOk bool
	}{
		{name: Upper, args: args{name: Upper}, wantF: toUpperFunc, wantOk: true},
		{name: Lower, args: args{name: Lower}, wantF: toLowerFunc, wantOk: true},
		{name: Snake, args: args{name: Snake}, wantF: toSnakeFunc, wantOk: true},
		{name: Camel, args: args{name: Camel}, wantF: toCamelFunc, wantOk: true},
		{name: "not_exist", args: args{name: "not_exist"}, wantF: nil, wantOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotF, gotOk := GetConvertFunc(tt.args.name)
			a := assert.New(t)
			a.EqualValues(gotOk, tt.wantOk)
			gotFName := runtime.FuncForPC(reflect.ValueOf(gotF).Pointer()).Name()
			wantFName := runtime.FuncForPC(reflect.ValueOf(tt.wantF).Pointer()).Name()
			a.Equal(gotFName, wantFName)
		})
	}
}

func TestRegisterConvertFunc(t *testing.T) {
	// init test data
	type Input struct {
		convertName string
		convertStr  string
		convertFunc func(string) string
	}
	data := []struct {
		input  Input
		output string
	}{
		{input: Input{"reg1", "a", func(s string) string { return s + "1" }}, output: "a1"},
		{input: Input{"reg2", "a", func(s string) string { return s + "2" }}, output: "a2"},
		{input: Input{"reg3", "a", func(s string) string { return s }}, output: "a"},
	}

	for _, tt := range data {
		t.Run(tt.input.convertName, func(t *testing.T) {
			RegisterConvertFunc(tt.input.convertName, tt.input.convertFunc)
			f, ok := GetConvertFunc(tt.input.convertName)
			if !ok {
				t.Errorf("register func error")
			}
			result := f(tt.input.convertStr)
			if result != tt.output {
				t.Errorf("expected %s, got %s", tt.output, result)
			}
		})
	}
}

func TestRegisterRepeatConvertFunc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	RegisterConvertFunc("reg1", func(s string) string { return s })
	RegisterConvertFunc("reg1", func(s string) string { return s })
}

func TestConvertKey_Array(t *testing.T) {
	data := []struct {
		input  string
		output string
	}{
		{input: `[1,2,3]`, output: `[1,2,3]`},
		{input: `["a","b","c"]`, output: `["a","b","c"]`},
		{input: `["{a","b}","{c}","{{d}}","{{e}"]`, output: `["{a","b}","{c}","{{d}}","{{e}"]`},
		{input: `[{"name":"Bob"},{"name":"Sam"}]`, output: `[{"NAME":"Bob"},{"NAME":"Sam"}]`},
		{input: `[{"name":"Bob"},{"scores":[{"English":"80"},{"Physics":75}]}]`,
			output: `[{"NAME":"Bob"},{"SCORES":[{"ENGLISH":"80"},{"PHYSICS":75}]}]`},
	}

	for _, tt := range data {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ConvertKey([]byte(tt.input), Upper)
			if err != nil {
				t.Errorf("convert error:%s", err.Error())
			}
			if string(result) != tt.output {
				t.Errorf("expected %s, got %s", tt.output, string(result))
			}
		})
	}
}

func TestConvertKey_Object(t *testing.T) {
	data := []struct {
		input  string
		output string
	}{
		{input: `{"name": "Bob","age": 30,"cars": ["Ford", "BMW"],"phone": 12345656778234223498}`,
			output: `{"NAME": "Bob","AGE": 30,"CARS": ["Ford", "BMW"],"PHONE": 12345656778234223498}`,
		},
		{input: `{"name": "Bob","address": {"city": "Beijing","street": "Haidian","detail": {"floor": 1,"room": 101}}}`,
			output: `{"NAME": "Bob","ADDRESS": {"CITY": "Beijing","STREET": "Haidian","DETAIL": {"FLOOR": 1,"ROOM": 101}}}`,
		},
	}

	for _, tt := range data {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ConvertKey([]byte(tt.input), Upper)
			if err != nil {
				t.Errorf("convert error:%s", err.Error())
			}
			a := assert.New(t)
			a.JSONEq(tt.output, string(result))
		})
	}

}

func TestConvertKey_ArrayObject(t *testing.T) {
	data := []struct {
		input  string
		output string
	}{
		{input: `{
			"name": "Bob",
			"contract":[
				{"relationship":"father","name":"Bob's father"},
				{"relationship":"wife","name":"[Bob's wife{1}]"}
			]}`,
			output: `{
				"NAME": "Bob",
				"CONTRACT":[
					{"RELATIONSHIP":"father","NAME":"Bob's father"},
					{"RELATIONSHIP":"wife","NAME":"[Bob's wife{1}]"}
				]}`,
		},
	}

	for _, tt := range data {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ConvertKey([]byte(tt.input), Upper)
			if err != nil {
				t.Errorf("convert error:%s", err.Error())
			}
			a := assert.New(t)
			a.JSONEq(tt.output, string(result))
		})
	}
}
