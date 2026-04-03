package ddlstructdiff

import (
	"go/ast"
	"go/token"
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
	strict  bool
)

func init() {
	Analyzer.Flags.StringVar(&ddlPath, "ddl", "", "ddl file path")
	Analyzer.Flags.BoolVar(&strict, "strict", false, "enable strict case sensitivity")
}

func hasIgnoreComment(cg *ast.CommentGroup) bool {
	if cg == nil {
		return false
	}
	for _, c := range cg.List {
		if strings.Contains(c.Text, "nolint:ddlstructdiff") {
			return true
		}
	}
	return false
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

	ddl, err := loadDDL(ddlPath, strict)
	if err != nil {
		return nil, err
	}

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	structs := NewStructs()
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return
		}

		structIgnored := hasIgnoreComment(genDecl.Doc)

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			st := NewStruct(typeSpec.Name.Name, typeSpec.Pos(), strict)
			if structIgnored {
				st.SetIgnored(true)
			}

			for _, field := range structType.Fields.List {
				tag := spannerTag(field)
				if tag != "" && len(field.Names) != 1 {
					pass.Reportf(field.Pos(), "field with spanner tag must have only one name")
					continue
				}
				fieldIgnored := hasIgnoreComment(field.Comment)
				for _, name := range field.Names {
					n := name.Name
					if tag != "" {
						n = tag
					}
					f := NewField(n, strict)
					if fieldIgnored {
						f.SetIgnored(true)
					}
					st.AddField(f)
				}
			}
			structs.AddStruct(st)
		}
	})

	for _, table := range ddl.Tables() {
		if table.IsIgnored() {
			continue
		}

		st, ok := structs.Struct(table.Name())
		if !ok {
			// TODO set option
			// pass.Reportf(token.NoPos, "%s struct corresponding to %s table not found", tableName, tableName)
			continue
		}

		if st.IsIgnored() {
			continue
		}

		for _, column := range table.Columns() {
			if column.IsIgnored() {
				continue
			}
			if _, ok := st.Field(column.Name()); !ok {
				pass.Reportf(st.Pos(), "%s struct must contain %s field corresponding to DDL", table.OriginalName(), column.OriginalName())
			}
		}
		for _, field := range st.Fields() {
			if field.IsIgnored() {
				continue
			}
			if _, ok := table.Column(field.Name()); !ok {
				pass.Reportf(st.Pos(), "%s table does not have a column corresponding to %s", table.OriginalName(), field.OriginalName())
			}
		}
	}

	return nil, nil
}
