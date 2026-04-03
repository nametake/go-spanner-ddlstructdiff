package ddlstructdiff

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
)

type Column struct {
	name    string
	strict  bool
	ignored bool
}

func NewColumn(name string, strict bool) *Column {
	return &Column{
		name:   name,
		strict: strict,
	}
}

func (c *Column) SetIgnored(ignored bool) {
	c.ignored = ignored
}

func (c *Column) IsIgnored() bool {
	return c.ignored
}

func (c *Column) Name() string {
	if c.strict {
		return c.name
	}
	return strings.ToLower(c.name)
}

func (c *Column) OriginalName() string {
	return c.name
}

type Table struct {
	name    string
	strict  bool
	ignored bool
	s       []*Column
	m       map[string]*Column
}

func NewTable(name string, strict bool) *Table {
	return &Table{
		name:   name,
		strict: strict,
		s:      []*Column{},
		m:      map[string]*Column{},
	}
}

func (t *Table) SetIgnored(ignored bool) {
	t.ignored = ignored
}

func (t *Table) IsIgnored() bool {
	return t.ignored
}

func (t *Table) Name() string {
	if t.strict {
		return t.name
	}
	return strings.ToLower(t.name)
}

func (t *Table) OriginalName() string {
	return t.name
}

func (t *Table) Columns() []*Column {
	return t.s
}

func (t *Table) Column(column string) (*Column, bool) {
	c, ok := t.m[column]
	return c, ok
}

func (t *Table) AddColumn(c *Column) {
	t.s = append(t.s, c)
	t.m[c.Name()] = c
}

type DDL struct {
	s []*Table
	m map[string]*Table
}

func NewDDL() *DDL {
	return &DDL{
		s: []*Table{},
		m: map[string]*Table{},
	}
}

func (d *DDL) Table(table string) (*Table, bool) {
	t, ok := d.m[table]
	return t, ok
}

func (d *DDL) Tables() []*Table {
	return d.s
}

func (d *DDL) AddTable(t *Table) {
	d.s = append(d.s, t)
	d.m[t.Name()] = t
}

func buildCommentMap(file *token.File) (map[token.Pos][]token.TokenComment, error) {
	lexer := &memefish.Lexer{File: file}
	commentMap := make(map[token.Pos][]token.TokenComment)
	for {
		if err := lexer.NextToken(); err != nil {
			return nil, err
		}
		if len(lexer.Token.Comments) > 0 {
			commentMap[lexer.Token.Pos] = lexer.Token.Comments
		}
		if lexer.Token.Kind == token.TokenEOF {
			break
		}
	}
	return commentMap, nil
}

func hasDDLIgnoreComment(commentMap map[token.Pos][]token.TokenComment, pos token.Pos) bool {
	comments, ok := commentMap[pos]
	if !ok {
		return false
	}
	for _, c := range comments {
		if strings.Contains(c.Raw, "nolint:ddlstructdiff") {
			return true
		}
	}
	return false
}

func loadDDL(ddlPath string, strict bool) (*DDL, error) {
	ddlFile, err := os.Open(ddlPath)
	if err != nil {
		return nil, err
	}

	ddlBytes, err := io.ReadAll(ddlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL file: %w", err)
	}

	buf := string(ddlBytes)

	commentMap, err := buildCommentMap(&token.File{Buffer: buf, FilePath: ddlPath})
	if err != nil {
		return nil, fmt.Errorf("failed to lex DDL: %w", err)
	}

	p := memefish.Parser{
		Lexer: &memefish.Lexer{File: &token.File{Buffer: buf, FilePath: ddlPath}},
	}

	stmt, err := p.ParseDDLs()
	if err != nil {
		return nil, fmt.Errorf("failed to parse DDL: %w", err)
	}

	ddl := NewDDL()
	for _, s := range stmt {
		ct, ok := s.(*ast.CreateTable)
		if !ok {
			continue
		}
		table := NewTable(ct.Name.SQL(), strict)
		if hasDDLIgnoreComment(commentMap, ct.Create) {
			table.SetIgnored(true)
		}
		for _, c := range ct.Columns {
			col := NewColumn(c.Name.Name, strict)
			if hasDDLIgnoreComment(commentMap, c.Name.NamePos) {
				col.SetIgnored(true)
			}
			table.AddColumn(col)
		}
		ddl.AddTable(table)
	}

	return ddl, nil
}
