package main

import (
	"github.com/nametake/ddlstructdiff"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(ddlstructdiff.Analyzer) }
