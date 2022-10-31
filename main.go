package main

import (
	"github.com/kimhongji/nomadcoin/cli"
	"github.com/kimhongji/nomadcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}