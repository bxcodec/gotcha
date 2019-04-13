# gotcha

gotcha: inmemory-cache in Go (Golang) with customizable algorithm

[![GoDoc](https://godoc.org/github.com/bxcodec/gotcha?status.svg)](https://godoc.org/github.com/bxcodec/gotcha)

## Index

* [Support](#support)
* [Getting Started](#getting-started)
* [Example](#example)
* [Contribution](#contribution)


## Support

You can file an [Issue](https://github.com/bxcodec/gotcha/issues/new).
See documentation in [Godoc](https://godoc.org/github.com/bxcodec/gotcha)


## Getting Started

#### Download

```shell
go get -u github.com/bxcodec/gotcha
```
## Example


### With Cache Client
```go
package main

import (
	"fmt"
	"log"

	"github.com/bxcodec/gotcha"
)

func main() {
	cache := gotcha.New()
	err := cache.Set("name", "John Snow")
	if err != nil {
		log.Fatal(err)
	}
	val, err := cache.Get("name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
```

### Without Cache Client
```go
package main

import (
	"fmt"
	"log"

	"github.com/bxcodec/gotcha"
)

func main() {
	err := gotcha.Set("name", "John Snow")
	if err != nil {
		log.Fatal(err)
	}
	val, err := gotcha.Get("name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
```


## Contribution
- You can submit an issue or create a Pull Request (PR)