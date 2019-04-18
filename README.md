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

### With Custom Cache ALgorithm

You can also custom and change the algorithm, expiry-time and also maximum memory.

```go
gotcha.NewOption().SetAlgorithm(cache.LRUAlgorithm).
	  SetExpiryTime(time.Minute * 10).
	  SetMaxSizeItem(100).
	  SetMaxMemory(cache.MB * 10)
```

**Warn**: Even gotcha support for MaxMemory, but the current version it's still using a simple json/encoding to count the byte size. So it will be slower if you set the MaxMemory.

Benchmark for LRU with/without MaxMemory

```
# With MaxMemory
20000000	      7878 ns/op	    1646 B/op	      20 allocs/op

# Without MaxMemory
200000000	       776 ns/op	     150 B/op	       6 allocs/op
```
If you seeking for fast performances and also your memory is high, ignore the MaxMemory options. I'm still looking for the better solutions for this problem. And if you have a better solutions, please kindly open and issue or submit a PR directly for the better results. 

#### LRU
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bxcodec/gotcha"
	"github.com/bxcodec/gotcha/cache"
)

func main() {
	cache := gotcha.New(
		gotcha.NewOption().SetAlgorithm(cache.LRUAlgorithm).
			SetExpiryTime(time.Minute * 10).SetMaxSizeItem(100),
	)
	err := cache.Set("Kue", "Nama")
	if err != nil {
		log.Fatal(err)
	}
	val, err := cache.Get("Kue")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
```

#### LFU
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bxcodec/gotcha"
	"github.com/bxcodec/gotcha/cache"
)

func main() {
	cache := gotcha.New(
		gotcha.NewOption().SetAlgorithm(cache.LFUAlgorithm).
			SetExpiryTime(time.Minute * 10).SetMaxSizeItem(100),
	)
	err := cache.Set("Kue", "Nama")
	if err != nil {
		log.Fatal(err)
	}
	val, err := cache.Get("Kue")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
}
```



## Contribution
- You can submit an issue or create a Pull Request (PR)
