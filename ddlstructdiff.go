package ddlstructdiff

import (
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

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	structs := NewStructs()
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return
		}

		st := NewStruct(typeSpec.Pos())
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				st.AddField(name.Name, NewField())
			}
		}

		structs.AddStruct(typeSpec.Name.Name, st)
	})

	for table, columns := range ddl {
		s, ok := structs.Struct(table)
		if !ok {
			continue
		}
		for column := range columns {
			_, ok := s.Field(column)
			if !ok {
				pass.Reportf(s.Pos, "%s struct must contain %s field corresponding to DDL", table, column)
			}
		}
	}

	return nil, nil
}
