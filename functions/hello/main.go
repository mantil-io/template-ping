package main

import (
	"github.com/atoz-technology/go-mantil-template/api/hello"
	"github.com/atoz-technology/mantil.go"
)

func main() {
	var api = hello.New()
	mantil.LambdaHandler(api)
}
