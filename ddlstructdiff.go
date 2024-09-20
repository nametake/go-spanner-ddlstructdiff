package ddlstructdiff

import (
	"go/ast"
	"os"
	"strings"

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

func spannerTag(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	tag := field.Tag.Value
	tag = strings.Trim(tag, "`")
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, `spanner:"`) {
			return strings.Trim(part[len(`spanner:"`):], `"`)
		}
	}
	return ""
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
			tag := spannerTag(field)
			if tag != "" && len(field.Names) != 1 {
				pass.Reportf(field.Pos(), "field with spanner tag must have only one name")
				continue
			}
			for _, name := range field.Names {
				n := name.Name
				if tag != "" {
					n = tag
				}
				st.AddField(NewField(n))
			}
		}

		structs.AddStruct(typeSpec.Name.Name, st)
	})

	for tableName, table := range ddl {
		s, ok := structs.Struct(tableName)
		if !ok {
			// TODO set option
			// pass.Reportf(token.NoPos, "%s struct corresponding to %s table not found", tableName, tableName)
			continue
		}
		for _, column := range table.Columns() {
			if _, ok := s.Field(column.LowerName()); !ok {
				pass.Reportf(s.Pos, "%s struct must contain %s field corresponding to DDL", tableName, column.Name)
			}
		}
		for _, field := range s.Fields() {
			if _, ok := table.Column(field.LowerName()); !ok {
				pass.Reportf(s.Pos, "%s table does not have a column corresponding to %s", tableName, field.Name)
			}
		}
	}

	return nil, nil
}
