package main

import (
	ddlstructdiff "github.com/nametake/go-spanner-ddlstructdiff"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(ddlstructdiff.Analyzer) }
