# json-key-convertor
![workflow](https://github.com/chengshicheng/json-key-convertor/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/chengshicheng/json-key-convertor)](https://goreportcard.com/report/github.com/chengshicheng/json-key-convertor)
[![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/chengshicheng/json-key-convertor)
[![codecov](https://codecov.io/gh/chengshicheng/json-key-convertor/branch/main/graph/badge.svg?token=SIOQNMM823)](https://codecov.io/gh/chengshicheng/json-key-convertor)

json-key-convertor is a Go package, it provides a simple way to converts the format of all keys in dynamic JSON to another format, such as case and underscore


# Getting Started

## Features

- Four built-in key conversion methods, lowercase, uppercase, snakecase, camlecase
- Support for nested multi-level json
- Custom key convert methods

## Installing

To start using json-key-convertor, install Go and run `go get`:

```sh
$ go get -u github.com/chengshicheng/json-key-convertor
```

This will retrieve the library.

## Example

### General Key Convert
```Go
package main

import (
	convertor "github.com/chengshicheng/json-key-convertor"
)

const jsonStr = `{"name":{"first_name":"Janet","last_name":"Prichard"},"age":47}`

func main() {
	value, err := convertor.ConvertKey([]byte(jsonStr), convertor.Camel)
	if err != nil {
		println(err.Error())
	}
	println(string(value))
	// {"Age":47,"Name":{"FirstName":"Janet","LastName":"Prichard"}}
}
```

### Custom Key Convert
```Go
package main

import (
	convertor "github.com/chengshicheng/json-key-convertor"
)

const jsonStr = `{"name":{"first_name":"Janet","last_name":"Prichard"},"age":47}`

func myPrefixFunc(s string) string {
	return "my_" + s
}

func main() {
	// register convert function
	convertor.RegisterConvertFunc("myprefix", myPrefixFunc)
	value, err := convertor.ConvertKey([]byte(jsonStr), "myprefix")
	if err != nil {
		println(err.Error())
	}
	println(string(value))
	// {"my_age":47,"my_name":{"my_first_name":"Janet","my_last_name":"Prichard"}}
}
```
