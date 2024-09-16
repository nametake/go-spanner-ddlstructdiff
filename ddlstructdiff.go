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
		(*ast.TypeSpec)(nil),
	}

	structs := NewEmptyStructs()
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return
		}

		st := NewEmptyStruct()
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				st.AddField(name.Name, &Field{})
			}
		}

		structs.AddStruct(typeSpec.Name.Name, st)

		fmt.Printf("%#v", structs)
		// table, ok := ddl.Table(typeSpec.Name.Name)
		// if !ok {
		// 	return
		// }
		//
		// fmt.Println(structType)
		// fmt.Println(table)
	})

	return nil, nil
}
