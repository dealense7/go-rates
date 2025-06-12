package main

import (
	_ "github.com/dealense7/go-rate-app/cmd/parser/currency"
	_ "github.com/dealense7/go-rate-app/cmd/parser/gas"
	"github.com/dealense7/go-rate-app/cmd/parser/root"
	_ "github.com/dealense7/go-rate-app/cmd/parser/store"
)

func main() {
	root.Execute()
}
