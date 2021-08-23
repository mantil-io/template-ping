package main

import (
	"github.com/mantil-io/go-mantil-template/api/hello"
	"github.com/mantil-io/mantil.go"
)

func main() {
	var api = hello.New()
	mantil.LambdaHandler(api)
}
