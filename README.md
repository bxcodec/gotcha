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
go get -u github.com/bxcodec/faker/v3
```
# Example

```go
cache:=gotcha.New()
cache.Set("key", 20)
res,err:=cache.Get("key")
if err != nil {
    log.Fatal(err)
}
fmt.Println(res)
```

## Contribution
- You can submit an issue or create a Pull Request (PR)