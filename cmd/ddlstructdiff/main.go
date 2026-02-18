package main

import (
	ddlstructdiff "github.com/nametake/go-spanner-ddlstructdiff"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(ddlstructdiff.Analyzer) }
