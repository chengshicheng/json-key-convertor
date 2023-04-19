# json-key-convertor
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

func main() {
	const jsonStr = `{"name":{"first_name":"Janet","last_name":"Prichard"},"age":47}`
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

func myPrefixFunc(s string) string {
	return "my_" + s
}

func main() {
	// register convert function
	convertor.RegisterConvertFunc("myprefix", myPrefixFunc)
	const jsonStr = `{"name":{"first_name":"Janet","last_name":"Prichard"},"age":47}`
	value, err := convertor.ConvertKey([]byte(jsonStr), "myprefix")
	if err != nil {
		println(err.Error())
	}
	println(string(value))
	// {"my_age":47,"my_name":{"my_first_name":"Janet","my_last_name":"Prichard"}}
}
```
