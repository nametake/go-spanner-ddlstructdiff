package ddlstructdiff

import (
	"fmt"
	"go/ast"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "ddlstructdiff is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "ddlstructdiff",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

var (
	ddlPath string
)

func init() {
	Analyzer.Flags.StringVar(&ddlPath, "ddl", "", "ddl file path")
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	ddlFile, err := os.Open(ddlPath)
	if err != nil {
		return nil, err
	}

	ddl, err := loadDDL(ddlFile)
	if err != nil {
		return nil, err
	}

	fmt.Println(ddl)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == "gopher" {
				pass.Reportf(n.Pos(), "identifier is gopher")
			}
		}
	})

	return nil, nil
}
