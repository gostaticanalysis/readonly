package main

import (
	"github.com/tenntenn/readonly"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(readonly.Analyzer) }
