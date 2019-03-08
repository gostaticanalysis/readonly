package main

import (
	"github.com/gostaticanalysis/readonly"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(readonly.Analyzer) }
