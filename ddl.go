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
	name   string
	strict bool
}

func NewColumn(name string, strict bool) *Column {
	return &Column{
		name:   name,
		strict: strict,
	}
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
	name   string
	strict bool
	s      []*Column
	m      map[string]*Column
}

func NewTable(name string, strict bool) *Table {
	return &Table{
		name:   name,
		strict: strict,
		s:      []*Column{},
		m:      map[string]*Column{},
	}
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

func loadDDL(ddlPath string, strict bool) (*DDL, error) {
	ddlFile, err := os.Open(ddlPath)
	if err != nil {
		return nil, err
	}

	ddlReader, err := io.ReadAll(ddlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL file: %w", err)
	}

	file := &token.File{
		Buffer:   string(ddlReader),
		FilePath: ddlPath,
	}

	p := memefish.Parser{
		Lexer: &memefish.Lexer{File: file},
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
		table := NewTable(ct.Name.Name, strict)
		for _, c := range ct.Columns {
			table.AddColumn(NewColumn(c.Name.Name, strict))
		}
		ddl.AddTable(table)
	}

	return ddl, nil
}
